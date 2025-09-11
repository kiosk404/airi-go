/**
 * Copyright (c) 2017, Smooth Code
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 * modified from https://github.com/gregberge/svgr/blob/7595d378b73d4826a4cead165b3f32386b07315b/packages/webpack/src/index.ts
 */
import type { Rspack } from '@rsbuild/core';
import { type Config } from '@svgr/core';
declare function svgrLoader(this: Rspack.LoaderContext<Config>, contents: string): void;
export default svgrLoader;
