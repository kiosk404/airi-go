"use strict";
let __rslib_import_meta_url__ = 'undefined' == typeof document ? new (require('url'.replace('', ''))).URL('file:' + __filename).href : document.currentScript && document.currentScript.src || new URL('main.js', document.baseURI).href;
var __webpack_require__ = {};
__webpack_require__.n = function(module) {
    var getter = module && module.__esModule ? function() {
        return module.default;
    } : function() {
        return module;
    };
    return __webpack_require__.d(getter, {
        a: getter
    }), getter;
}, __webpack_require__.d = function(exports1, definition) {
    for(var key in definition)__webpack_require__.o(definition, key) && !__webpack_require__.o(exports1, key) && Object.defineProperty(exports1, key, {
        enumerable: !0,
        get: definition[key]
    });
}, __webpack_require__.o = function(obj, prop) {
    return Object.prototype.hasOwnProperty.call(obj, prop);
}, __webpack_require__.r = function(exports1) {
    'undefined' != typeof Symbol && Symbol.toStringTag && Object.defineProperty(exports1, Symbol.toStringTag, {
        value: 'Module'
    }), Object.defineProperty(exports1, '__esModule', {
        value: !0
    });
};
var __webpack_exports__ = {};
__webpack_require__.r(__webpack_exports__), __webpack_require__.d(__webpack_exports__, {
    pluginSass: ()=>pluginSass,
    PLUGIN_SASS_NAME: ()=>PLUGIN_SASS_NAME
});
let external_node_module_namespaceObject = require("node:module"), external_node_path_namespaceObject = require("node:path");
var external_node_path_default = __webpack_require__.n(external_node_path_namespaceObject);
let external_node_url_namespaceObject = require("node:url"), external_deepmerge_namespaceObject = require("deepmerge");
var external_deepmerge_default = __webpack_require__.n(external_deepmerge_namespaceObject);
let external_reduce_configs_namespaceObject = require("reduce-configs"), index_js_namespaceObject = require("../compiled/resolve-url-loader/index.js");
var index_js_default = __webpack_require__.n(index_js_namespaceObject);
let GLOBAL_PATCHED_SYMBOL = Symbol('GLOBAL_PATCHED_SYMBOL');
function patchGlobalLocation() {
    if (!global.location) {
        let location = Object.freeze({
            [GLOBAL_PATCHED_SYMBOL]: !0,
            href: (0, external_node_url_namespaceObject.pathToFileURL)(process.cwd()).href + external_node_path_default().sep
        });
        global.location = location;
    }
}
function unpatchGlobalLocation() {
    var _global_location;
    (null === (_global_location = global.location) || void 0 === _global_location ? void 0 : _global_location[GLOBAL_PATCHED_SYMBOL]) && delete global.location;
}
function patchCompilerGlobalLocation(compiler) {
    compiler.hooks.run.tap('PatchGlobalLocation', patchGlobalLocation), compiler.hooks.watchRun.tap('PatchGlobalLocation', patchGlobalLocation), compiler.hooks.watchClose.tap('PatchGlobalLocation', unpatchGlobalLocation), compiler.hooks.done.tap('PatchGlobalLocation', unpatchGlobalLocation);
}
let getResolveUrlJoinFn = async ()=>{
    let { createJoinFunction, asGenerator, createJoinImplementation, defaultJoinGenerator } = index_js_default();
    return createJoinFunction('rsbuild-resolve-join-fn', createJoinImplementation(asGenerator((item, ...rest)=>item.uri.startsWith('.') ? defaultJoinGenerator(item, ...rest) : [
            null
        ])));
}, src_dirname = external_node_path_default().dirname((0, external_node_url_namespaceObject.fileURLToPath)(__rslib_import_meta_url__)), src_require = (0, external_node_module_namespaceObject.createRequire)(__rslib_import_meta_url__), PLUGIN_SASS_NAME = 'rsbuild:sass', getSassLoaderOptions = (userOptions, isUseCssSourceMap)=>{
    let excludes = [], mergedOptions = (0, external_reduce_configs_namespaceObject.reduceConfigsWithContext)({
        initial: {
            sourceMap: isUseCssSourceMap,
            api: 'modern-compiler',
            implementation: src_require.resolve('sass-embedded'),
            sassOptions: {
                quietDeps: !0
            }
        },
        config: userOptions,
        ctx: {
            addExcludes: (items)=>{
                excludes.push(...Array.isArray(items) ? items : [
                    items
                ]);
            }
        },
        mergeFn: (defaults, userOptions)=>({
                ...defaults,
                ...userOptions,
                sassOptions: defaults.sassOptions && userOptions.sassOptions ? external_deepmerge_default()(defaults.sassOptions, userOptions.sassOptions) : userOptions.sassOptions || defaults.sassOptions
            })
    });
    return mergedOptions.sassOptions ||= {}, !mergedOptions.sassOptions.silenceDeprecations && (mergedOptions.sassOptions.silenceDeprecations = [
        'import'
    ], 'legacy' === mergedOptions.api && mergedOptions.sassOptions.silenceDeprecations.push('legacy-js-api')), {
        options: mergedOptions,
        excludes
    };
}, findRuleId = (chain, defaultId)=>{
    let id = defaultId, index = 0;
    for(; chain.module.rules.has(id);)id = `${defaultId}-${++index}`;
    return id;
}, pluginSass = (pluginOptions = {})=>({
        name: PLUGIN_SASS_NAME,
        setup (api) {
            api.onAfterCreateCompiler(({ compiler })=>{
                patchCompilerGlobalLocation(compiler);
            }), api.modifyBundlerChain(async (chain, { CHAIN_ID })=>{
                let { excludes, options } = getSassLoaderOptions(pluginOptions.sassLoaderOptions, !0), ruleId = findRuleId(chain, CHAIN_ID.RULE.SASS), rule = chain.module.rule(ruleId).test(pluginOptions.include ?? /\.s(?:a|c)ss$/).merge({
                    sideEffects: !0
                }).resolve.preferRelative(!0).end();
                for (let item of excludes)rule.exclude.add(item);
                pluginOptions.exclude && rule.exclude.add(pluginOptions.exclude);
                let cssRule = chain.module.rules.get(CHAIN_ID.RULE.CSS);
                for (let id of (rule.dependency(cssRule.get('dependency')), Object.keys(cssRule.uses.entries()))){
                    let loader = cssRule.uses.get(id), options = loader.get('options') ?? {}, clonedOptions = external_deepmerge_default()({}, options);
                    id === CHAIN_ID.USE.CSS && (clonedOptions.importLoaders += 2), rule.use(id).loader(loader.get('loader')).options(clonedOptions);
                }
                rule.use(CHAIN_ID.USE.RESOLVE_URL).loader(external_node_path_default().join(src_dirname, '../compiled/resolve-url-loader/index.js')).options({
                    join: await getResolveUrlJoinFn(),
                    sourceMap: !1
                }).end().use(CHAIN_ID.USE.SASS).loader(external_node_path_default().join(src_dirname, '../compiled/sass-loader/index.js')).options(options);
            });
        }
    });
var __webpack_export_target__ = exports;
for(var __webpack_i__ in __webpack_exports__)__webpack_export_target__[__webpack_i__] = __webpack_exports__[__webpack_i__];
__webpack_exports__.__esModule && Object.defineProperty(__webpack_export_target__, '__esModule', {
    value: !0
});
