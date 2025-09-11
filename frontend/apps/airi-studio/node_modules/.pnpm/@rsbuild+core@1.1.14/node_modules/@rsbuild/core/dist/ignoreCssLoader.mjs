function ignoreCssLoader_rslib_entry_(source) {
    return (null === this || void 0 === this || this.cacheable(!0), source.includes('___CSS_LOADER_EXPORT___')) ? '' : source;
}
export { ignoreCssLoader_rslib_entry_ as default };
