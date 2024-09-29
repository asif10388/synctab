import fs from "fs";
import path from "path";
import webpack from "webpack";
import CopyWebpackPlugin from "copy-webpack-plugin";

const getAllScripts = () => {
  const fileMap: {
    [key: string]: string;
  } = {};

  const scriptPath = "./src/scripts";
  const files = fs.readdirSync(scriptPath, { recursive: true }) as string[];

  files.forEach((file) => {
    const fileName = file.split("/")[file.split("/").length - 1].split(".")[0];
    if (file.split(".")[1] === "ts") fileMap[fileName] = `${scriptPath}/${file}`;
  });

  return fileMap;
};

const config: webpack.Configuration = {
  entry: getAllScripts(),

  resolve: {
    extensions: [".ts", ".js"],
  },

  module: {
    rules: [
      {
        test: /\.ts$/,
        loader: "ts-loader",
        exclude: /node_modules/,
      },
    ],
  },

  output: {
    filename: "[name].js",
    path: path.resolve(__dirname, "dist"),
    clean: true, // Clean the output directory before emit.
  },

  plugins: [
    new CopyWebpackPlugin({
      patterns: [
        { from: "public" },
        { from: "./src/pages" },
        { from: "./src/styles" },
        { from: "manifest.json" },
      ],
    }),
  ],
};

export default config;
