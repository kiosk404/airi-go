#!/usr/bin/env node
/**
 * Post-process script to organize generated TypeScript files by namespace
 * Parses IDL files to build type->namespace mapping, then moves generated files accordingly
 */

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const GENERATED_DIR = path.join(__dirname, '../src/api/generated');
const IDL_DIR = path.join(__dirname, '../../../idl');

/**
 * Parse a thrift file and extract type names with their namespace
 */
function parseThriftFile(filePath, relativePath) {
    const content = fs.readFileSync(filePath, 'utf-8');
    const types = [];

    // Extract namespace (use 'go' namespace as reference)
    const namespaceMatch = content.match(/namespace\s+go\s+([\w.]+)/);
    if (!namespaceMatch) return types;

    const namespace = namespaceMatch[1];

    // Convert namespace to directory path (e.g., "app.developer_api" -> "app/developer_api")
    const dirPath = namespace.replace(/\./g, '/');

    // Extract struct names
    const structMatches = content.matchAll(/struct\s+(\w+)\s*\{/g);
    for (const match of structMatches) {
        types.push({ name: match[1], namespace, dirPath });
    }

    // Extract enum names
    const enumMatches = content.matchAll(/enum\s+(\w+)\s*\{/g);
    for (const match of enumMatches) {
        types.push({ name: match[1], namespace, dirPath });
    }

    // Extract service names
    const serviceMatches = content.matchAll(/service\s+(\w+)\s*\{/g);
    for (const match of serviceMatches) {
        types.push({ name: match[1], namespace, dirPath });
    }

    // Extract typedef names
    const typedefMatches = content.matchAll(/typedef\s+[\w<>,\s]+\s+(\w+)/g);
    for (const match of typedefMatches) {
        types.push({ name: match[1], namespace, dirPath });
    }

    // Extract exception names
    const exceptionMatches = content.matchAll(/exception\s+(\w+)\s*\{/g);
    for (const match of exceptionMatches) {
        types.push({ name: match[1], namespace, dirPath });
    }

    return types;
}

/**
 * Recursively scan IDL directory and build type->namespace mapping
 */
function buildTypeMapping() {
    const typeMap = new Map();

    function scanDir(dir, basePath = '') {
        const entries = fs.readdirSync(dir, { withFileTypes: true });

        for (const entry of entries) {
            const fullPath = path.join(dir, entry.name);
            const relativePath = path.join(basePath, entry.name);

            if (entry.isDirectory()) {
                scanDir(fullPath, relativePath);
            } else if (entry.name.endsWith('.thrift')) {
                const types = parseThriftFile(fullPath, relativePath);
                for (const type of types) {
                    typeMap.set(type.name, type);
                }
            }
        }
    }

    scanDir(IDL_DIR);
    return typeMap;
}

/**
 * Update imports in a TypeScript file after moving
 */
function updateImports(filePath, movedFiles) {
    let content = fs.readFileSync(filePath, 'utf-8');
    let modified = false;

    // Find all imports and update paths
    const importRegex = /from\s+["']\.\/(\w+)["']/g;
    content = content.replace(importRegex, (match, importName) => {
        const importInfo = movedFiles.get(importName);
        if (importInfo && importInfo.dirPath) {
            const currentFileInfo = movedFiles.get(path.basename(filePath, '.ts'));
            if (currentFileInfo && currentFileInfo.dirPath === importInfo.dirPath) {
                // Same directory, keep relative
                return match;
            } else if (currentFileInfo && currentFileInfo.dirPath) {
                // Different directory, calculate relative path
                const fromDir = currentFileInfo.dirPath;
                const toDir = importInfo.dirPath;
                const relativePath = path.relative(fromDir, toDir) || '.';
                modified = true;
                return `from "${relativePath}/${importName}"`;
            } else {
                // Current file is in root, import from subdirectory
                modified = true;
                return `from "./${importInfo.dirPath}/${importName}"`;
            }
        }
        return match;
    });

    if (modified) {
        fs.writeFileSync(filePath, content);
    }
}

function organizeFiles() {
    if (!fs.existsSync(GENERATED_DIR)) {
        console.log('Generated directory does not exist:', GENERATED_DIR);
        return;
    }

    console.log('Building type mapping from IDL files...');
    const typeMap = buildTypeMapping();
    console.log(`Found ${typeMap.size} types in IDL files`);

    const files = fs.readdirSync(GENERATED_DIR).filter(f => f.endsWith('.ts') && f !== 'index.ts');
    const movedFiles = new Map(); // filename (without .ts) -> { dirPath, newPath }
    const dirFiles = new Map(); // dirPath -> [filenames]
    const rootFiles = [];

    console.log(`\nProcessing ${files.length} TypeScript files...`);

    // First pass: determine where each file should go
    for (const file of files) {
        const typeName = file.replace('.ts', '');
        const typeInfo = typeMap.get(typeName);

        if (typeInfo) {
            movedFiles.set(typeName, { dirPath: typeInfo.dirPath, namespace: typeInfo.namespace });
            if (!dirFiles.has(typeInfo.dirPath)) {
                dirFiles.set(typeInfo.dirPath, []);
            }
            dirFiles.get(typeInfo.dirPath).push(file);
        } else {
            movedFiles.set(typeName, { dirPath: null, namespace: null });
            rootFiles.push(file);
        }
    }

    // Second pass: move files
    for (const [dirPath, fileList] of dirFiles) {
        const targetDir = path.join(GENERATED_DIR, dirPath);
        fs.mkdirSync(targetDir, { recursive: true });

        for (const file of fileList) {
            const srcPath = path.join(GENERATED_DIR, file);
            const destPath = path.join(targetDir, file);
            fs.renameSync(srcPath, destPath);
            console.log(`  ${file} -> ${dirPath}/`);
        }
    }

    // Third pass: update imports in moved files
    console.log('\nUpdating imports...');
    for (const [dirPath, fileList] of dirFiles) {
        const targetDir = path.join(GENERATED_DIR, dirPath);
        for (const file of fileList) {
            updateImports(path.join(targetDir, file), movedFiles);
        }
    }

    // Update imports in root files
    for (const file of rootFiles) {
        updateImports(path.join(GENERATED_DIR, file), movedFiles);
    }

    // Create index.ts for each namespace directory
    for (const [dirPath, fileList] of dirFiles) {
        const nsDir = path.join(GENERATED_DIR, dirPath);
        const exports = fileList.map(f => `export * from './${f.replace('.ts', '')}';`).sort().join('\n');
        fs.writeFileSync(path.join(nsDir, 'index.ts'), exports + '\n');
        console.log(`Created ${dirPath}/index.ts (${fileList.length} exports)`);
    }

    // Create root index.ts
    const rootExports = [];

    // Export from namespace directories
    for (const dirPath of [...dirFiles.keys()].sort()) {
        // Use the first part of the path as the export name
        const parts = dirPath.split('/');
        if (parts.length === 1) {
            rootExports.push(`export * from './${dirPath}';`);
        }
    }

    // Also create intermediate index files for nested namespaces
    const topLevelDirs = new Set();
    for (const dirPath of dirFiles.keys()) {
        const topLevel = dirPath.split('/')[0];
        topLevelDirs.add(topLevel);
    }

    for (const topDir of topLevelDirs) {
        const subDirs = [...dirFiles.keys()].filter(d => d.startsWith(topDir + '/'));
        if (subDirs.length > 0) {
            const topDirPath = path.join(GENERATED_DIR, topDir);
            if (!fs.existsSync(path.join(topDirPath, 'index.ts'))) {
                const subExports = subDirs.map(d => {
                    const subPath = d.substring(topDir.length + 1);
                    return `export * from './${subPath}';`;
                }).join('\n');
                fs.mkdirSync(topDirPath, { recursive: true });
                fs.writeFileSync(path.join(topDirPath, 'index.ts'), subExports + '\n');
                console.log(`Created ${topDir}/index.ts (aggregating subdirectories)`);
            }
        }
        rootExports.push(`export * from './${topDir}';`);
    }

    // Export root files
    for (const file of rootFiles.sort()) {
        rootExports.push(`export * from './${file.replace('.ts', '')}';`);
    }

    fs.writeFileSync(path.join(GENERATED_DIR, 'index.ts'), [...new Set(rootExports)].sort().join('\n') + '\n');

    console.log(`\n✅ Organization complete!`);
    console.log(`   - ${dirFiles.size} namespace directories created`);
    console.log(`   - ${files.length - rootFiles.length} files organized into namespaces`);
    console.log(`   - ${rootFiles.length} files kept in root (no namespace match)`);
}

organizeFiles();
