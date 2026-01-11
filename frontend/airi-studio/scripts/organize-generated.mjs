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
/**
 * Convert camelCase/PascalCase field names to snake_case.
 * NOTE: This is a simple converter and intentionally keeps consecutive capitals as "_i_d" etc.
 */
function convertToSnakeCase(str) {
    return str.replace(/[A-Z]/g, (match) => '_' + match.toLowerCase()).replace(/^_/, '');
}

/**
 * Transform interface field declarations to snake_case.
 * Examples:
 * - "SyncAll?:" -> "sync_all?:"
 * - "BaseResp:" -> "base_resp:"
 */
function transformInterfaceFieldNamesToSnakeCase(content) {
    // Only transform within "export interface ... { ... }" blocks.
    // This avoids breaking codec object literals like "LogID:" / "BaseResp:".
    const interfaceBlockRegex = /export interface\s+\w+(?:\s+extends[^\{]+)?\s*\{[\s\S]*?\n\}/g;

    return content.replace(interfaceBlockRegex, (block) => {
        return block.replace(/^\s+(\w+)(\??):\s/gm, (match, fieldName) => {
            if (fieldName === fieldName.toLowerCase()) {
                return match;
            }
            return match.replace(fieldName, convertToSnakeCase(fieldName));
        });
    });
}

/**
 * Transform args.<FieldName> accesses to args.<field_name> so they match the transformed Args interfaces.
 * Example: args.SyncAll -> args.sync_all
 */
function transformArgsPropertyAccessToSnakeCase(content) {
    return content.replace(/\bargs\.(\w+)\b/g, (match, fieldName) => {
        if (fieldName === fieldName.toLowerCase()) {
            return match;
        }
        return `args.${convertToSnakeCase(fieldName)}`;
    });
}

/**
 * Transform keys inside "return { ... }" object literals to snake_case.
 * This is safe because the returned object isn't accessed by PascalCase later in the same function.
 */
function transformReturnObjectKeysToSnakeCase(content) {
    return content.replace(/return\s+\{([\s\S]*?)\};/g, (fullMatch, body) => {
        const newBody = body.replace(/^\s+(\w+):\s/gm, (match, key) => {
            if (key === key.toLowerCase()) {
                return match;
            }
            return match.replace(key, convertToSnakeCase(key));
        });
        return `return {${newBody}};`;
    });
}

function transformClassPropertyDeclarationsToSnakeCase(content) {
    // Transform class fields: "public FieldName?:" -> "public field_name?:"
    return content.replace(/^\s+public\s+(\w+)\??:\s/gm, (match, fieldName) => {
        if (fieldName === fieldName.toLowerCase()) {
            return match;
        }
        return match.replace(fieldName, convertToSnakeCase(fieldName));
    });
}

function transformThisPropertyAccessToSnakeCase(content) {
    // Transform "this.FieldName" -> "this.field_name"
    return content.replace(/\bthis\.(\w+)\b/g, (match, fieldName) => {
        if (fieldName === fieldName.toLowerCase()) {
            return match;
        }
        return `this.${convertToSnakeCase(fieldName)}`;
    });
}

function transformFieldAnnotationsKeysToSnakeCase(content) {
    // Transform keys inside _fieldAnnotations blocks only.
    return content.replace(/(\b_fieldAnnotations\b[\s\S]*?=\s*\{)([\s\S]*?)(\};)/g, (fullMatch, start, body, end) => {
        const newBody = body.replace(/^\s+(\w+):\s*\{/gm, (match, key) => {
            if (key === key.toLowerCase()) {
                return match;
            }
            return match.replace(key, convertToSnakeCase(key));
        });
        return `${start}${newBody}${end}`;
    });
}

function transformFieldNamesToSnakeCase(content) {
    // Keep the original function name for compatibility with existing call sites.
    // Order matters: declare properties first, then update "this." accesses.
    content = transformInterfaceFieldNamesToSnakeCase(content);
    content = transformArgsPropertyAccessToSnakeCase(content);
    content = transformClassPropertyDeclarationsToSnakeCase(content);
    content = transformThisPropertyAccessToSnakeCase(content);
    content = transformFieldAnnotationsKeysToSnakeCase(content);
    content = transformReturnObjectKeysToSnakeCase(content);
    return content;
}

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

/**
 * Fix __ROOT_NAMESPACE__ imports that point to "./" but should point to other namespaces
 * This happens when thrift files include other thrift files - the generator incorrectly
 * resolves cross-namespace references to the current directory instead of the correct namespace
 */
function fixRootNamespaceImports(filePath, currentDirPath, typeMap) {
    let content = fs.readFileSync(filePath, 'utf-8');

    // Check if file imports from "./" as __ROOT_NAMESPACE__
    const hasRootNamespaceImport = content.includes('import * as __ROOT_NAMESPACE__ from "./"');
    if (!hasRootNamespaceImport) return false;

    // Find all types used with __ROOT_NAMESPACE__ prefix
    const usedTypes = new Set();
    const typeUsageRegex = /\b__ROOT_NAMESPACE__\.(\w+)/g;
    let match;
    while ((match = typeUsageRegex.exec(content)) !== null) {
        // Extract base type name (remove I prefix for interfaces, Codec/Args suffix)
        let typeName = match[1];
        if (typeName.startsWith('I') && typeName.length > 1 && typeName[1] === typeName[1].toUpperCase()) {
            // Could be an interface like IMsgParticipantInfo -> MsgParticipantInfo
            const baseName = typeName.substring(1);
            if (typeMap.has(baseName)) {
                usedTypes.add(baseName);
            } else if (typeMap.has(typeName)) {
                usedTypes.add(typeName);
            }
        } else if (typeName.endsWith('Codec')) {
            const baseName = typeName.replace(/Codec$/, '');
            if (typeMap.has(baseName)) {
                usedTypes.add(baseName);
            }
        } else if (typeName.endsWith('Args')) {
            const baseName = typeName.replace(/Args$/, '');
            if (typeMap.has(baseName)) {
                usedTypes.add(baseName);
            }
        } else if (typeMap.has(typeName)) {
            usedTypes.add(typeName);
        }
    }

    if (usedTypes.size === 0) return false;

    // Group types by their namespace directory
    const typesByDir = new Map();
    for (const typeName of usedTypes) {
        const typeInfo = typeMap.get(typeName);
        if (typeInfo && typeInfo.dirPath !== currentDirPath) {
            if (!typesByDir.has(typeInfo.dirPath)) {
                typesByDir.set(typeInfo.dirPath, []);
            }
            typesByDir.get(typeInfo.dirPath).push(typeName);
        }
    }

    if (typesByDir.size === 0) return false;

    // If all types are from a single directory, we can just change the import path
    if (typesByDir.size === 1) {
        const [targetDir] = typesByDir.keys();
        const relativePath = path.relative(currentDirPath, targetDir) || '.';
        const normalizedPath = relativePath.startsWith('.') ? relativePath : './' + relativePath;

        content = content.replace(
            'import * as __ROOT_NAMESPACE__ from "./"',
            `import * as __ROOT_NAMESPACE__ from "${normalizedPath}"`
        );

        fs.writeFileSync(filePath, content);
        console.log(`  Fixed __ROOT_NAMESPACE__ import in ${path.basename(filePath)} -> ${normalizedPath}`);
        return true;
    }

    // Multiple namespaces - need to add separate imports and replace usages
    // This is more complex, we need to:
    // 1. Remove the __ROOT_NAMESPACE__ import
    // 2. Add individual imports for each namespace
    // 3. Replace __ROOT_NAMESPACE__.Type with NamespaceName.Type

    const importLines = [];
    const replacements = new Map(); // old prefix -> new prefix

    for (const [targetDir, types] of typesByDir) {
        const relativePath = path.relative(currentDirPath, targetDir) || '.';
        const normalizedPath = relativePath.startsWith('.') ? relativePath : './' + relativePath;
        // Create a namespace alias from the directory path
        const nsAlias = targetDir.replace(/\//g, '_').toUpperCase() + '_NS';
        importLines.push(`import * as ${nsAlias} from "${normalizedPath}";`);

        for (const typeName of types) {
            // Map all variants of the type name
            replacements.set(`__ROOT_NAMESPACE__.${typeName}`, `${nsAlias}.${typeName}`);
            replacements.set(`__ROOT_NAMESPACE__.I${typeName}`, `${nsAlias}.I${typeName}`);
            replacements.set(`__ROOT_NAMESPACE__.${typeName}Codec`, `${nsAlias}.${typeName}Codec`);
            replacements.set(`__ROOT_NAMESPACE__.I${typeName}Args`, `${nsAlias}.I${typeName}Args`);
        }
    }

    // Remove the old import
    content = content.replace(
        /import \* as __ROOT_NAMESPACE__ from "\.\/";?\n?/,
        importLines.join('\n') + '\n'
    );

    // Apply replacements
    for (const [oldStr, newStr] of replacements) {
        content = content.split(oldStr).join(newStr);
    }

    fs.writeFileSync(filePath, content);
    console.log(`  Fixed multi-namespace imports in ${path.basename(filePath)}`);
    return true;
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

    // Second pass: move files and transform to snake_case
    console.log('\nTransforming field names to snake_case...');
    for (const [dirPath, fileList] of dirFiles) {
        const targetDir = path.join(GENERATED_DIR, dirPath);
        fs.mkdirSync(targetDir, { recursive: true });

        for (const file of fileList) {
            const srcPath = path.join(GENERATED_DIR, file);
            const destPath = path.join(targetDir, file);

            // Transform content before moving
            let content = fs.readFileSync(srcPath, 'utf-8');
            content = transformFieldNamesToSnakeCase(content);
            fs.writeFileSync(destPath, content);
            fs.unlinkSync(srcPath);

            console.log(`  ${file} -> ${dirPath}/`);
        }
    }

    // Also transform root files
    console.log('Transforming root files...');
    for (const file of rootFiles) {
        const filePath = path.join(GENERATED_DIR, file);
        let content = fs.readFileSync(filePath, 'utf-8');
        content = transformFieldNamesToSnakeCase(content);
        fs.writeFileSync(filePath, content);
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

    // Fourth pass: fix __ROOT_NAMESPACE__ imports (cross-namespace includes)
    console.log('\nFixing __ROOT_NAMESPACE__ imports...');
    let fixedCount = 0;
    for (const [dirPath, fileList] of dirFiles) {
        const targetDir = path.join(GENERATED_DIR, dirPath);
        for (const file of fileList) {
            if (fixRootNamespaceImports(path.join(targetDir, file), dirPath, typeMap)) {
                fixedCount++;
            }
        }
    }
    // Also check root files
    for (const file of rootFiles) {
        if (fixRootNamespaceImports(path.join(GENERATED_DIR, file), '', typeMap)) {
            fixedCount++;
        }
    }
    if (fixedCount > 0) {
        console.log(`  Fixed ${fixedCount} files with __ROOT_NAMESPACE__ imports`);
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
