import { merge } from "webpack-merge";
import config from "./webpack.cmn";
import { Configuration } from "webpack";

const merged = merge<Configuration>(config, {
  mode: "development",
  devtool: "inline-source-map",
});

export default merged;
