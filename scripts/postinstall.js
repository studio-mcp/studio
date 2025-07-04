#!/usr/bin/env node

// Adapted from https://github.com/supabase/cli/blob/develop/scripts/postinstall.js
// Copyright (c) 2021 Supabase, Inc. and contributors
// Released under the MIT License
// Adapted for @studio-mcp/studio 2025-07-02

// NOTE: we output to stderrr because stdout is reserved for MCP NDJSON.

// Ref 1: https://github.com/sanathkr/go-npm
// Ref 2: https://medium.com/xendit-engineering/how-we-repurposed-npm-to-publish-and-distribute-our-go-binaries-for-internal-cli-23981b80911b
"use strict";

import binLinks from "bin-links";
import { createHash } from "crypto";
import fs from "fs";
import fetch from "node-fetch";
import { Agent } from "https";
import { HttpsProxyAgent } from "https-proxy-agent";
import path from "path";
import { extract } from "tar";
import zlib from "zlib";

// Mapping from Node's `process.arch` to Golang's `$GOARCH`
const ARCH_MAPPING = {
  x64: "amd64",
  arm64: "arm64",
};

// Mapping between Node's `process.platform` to Golang's
const PLATFORM_MAPPING = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const arch = ARCH_MAPPING[process.arch];
const platform = PLATFORM_MAPPING[process.platform];

// Read package.json
const readPackageJson = async () => {
  const contents = await fs.promises.readFile("package.json");
  return JSON.parse(contents);
};

// Build the download url from package.json
const getDownloadUrl = (packageJson) => {
  const pkgName = packageJson.name;
  // Use the actual installed version from npm environment, fallback to package.json
  const version = process.env.npm_package_version || packageJson.version;
  const repo = packageJson.repository;
  const url = `https://github.com/${repo}/releases/download/v${version}/${pkgName}_${platform}_${arch}.tar.gz`;
  return url;
};

const fetchAndParseCheckSumFile = async (packageJson, agent) => {
  // Use the actual installed version from npm environment, fallback to package.json
  const version = process.env.npm_package_version || packageJson.version;
  const pkgName = packageJson.name;
  const repo = packageJson.repository;
  const checksumFileUrl = `https://github.com/${repo}/releases/download/v${version}/${pkgName}_${version}_checksums.txt`;

  // Fetch the checksum file
  console.error("Downloading", checksumFileUrl);
  const response = await fetch(checksumFileUrl, { agent });
  if (response.ok) {
    const checkSumContent = await response.text();
    const lines = checkSumContent.split("\n");

    const checksums = {};
    for (const line of lines) {
      const [checksum, packageName] = line.split(/\s+/);
      if (checksum && packageName) {
        checksums[packageName] = checksum;
      }
    }

    return checksums;
  } else {
    console.error(
      "Could not fetch checksum file",
      response.status,
      response.statusText
    );
  }
};

const errGlobal = `Installing Studio MCP CLI as a global module is not supported.
Please use one of the supported package managers or download directly from GitHub releases.
`;
const errChecksum = "Checksum mismatch. Downloaded data might be corrupted.";
const errUnsupported = `Installation is not supported for ${process.platform} ${process.arch}`;

/**
 * Reads the configuration from application's package.json,
 * downloads the binary from package url and stores at
 * ./bin in the package's root.
 *
 *  See: https://docs.npmjs.com/files/package.json#bin
 */
async function main() {
  const yarnGlobal = JSON.parse(
    process.env.npm_config_argv || "{}"
  ).original?.includes("global");

  // Allow npx usage - npx creates temporary local installs, not true global installs
  const isNpxContext = process.env.npm_command === "exec" ||
                      process.env.npm_execpath?.includes("npx") ||
                      process.env._?.includes("npx");

  if ((process.env.npm_config_global || yarnGlobal) && !isNpxContext) {
    throw errGlobal;
  }
  if (!arch || !platform) {
    throw errUnsupported;
  }

  // Read from package.json and prepare for the installation.
  const pkg = await readPackageJson();
  if (platform === "windows") {
    // Update bin path in package.json for Windows
    pkg.bin[Object.keys(pkg.bin)[0]] += ".exe";
  }

  // Prepare the installation path by creating the directory if it doesn't exist.
  const binPath = pkg.bin[Object.keys(pkg.bin)[0]];
  const binDir = path.dirname(binPath);
  await fs.promises.mkdir(binDir, { recursive: true });

  // Create the agent that will be used for all the fetch requests later.
  const proxyUrl =
    process.env.npm_config_https_proxy ||
    process.env.npm_config_http_proxy ||
    process.env.npm_config_proxy;
  // Keeps the TCP connection alive when sending multiple requests
  // Ref: https://github.com/node-fetch/node-fetch/issues/1735
  const agent = proxyUrl
    ? new HttpsProxyAgent(proxyUrl, { keepAlive: true })
    : new Agent({ keepAlive: true });

  // First, fetch the checksum map.
  const checksumMap = await fetchAndParseCheckSumFile(pkg, agent);

  // Then, download the binary.
  const url = getDownloadUrl(pkg);
  console.error("Downloading", url);
  const resp = await fetch(url, { agent });

  if (!resp.ok) {
    throw new Error(`Failed to download binary: ${resp.status} ${resp.statusText}`);
  }

  const hash = createHash("sha256");
  const pkgNameWithPlatform = `studio_${platform}_${arch}.tar.gz`;

  // Then, decompress the binary -- we will first Un-GZip, then we will untar.
  const ungz = zlib.createGunzip();
  const binName = path.basename(binPath);
  const untar = extract({ cwd: binDir }, [binName]);

  // Update the hash with the binary data as it's being downloaded.
  resp.body
    .on("data", (chunk) => {
      hash.update(chunk);
    })
    // Pipe the data to the ungz stream.
    .pipe(ungz);

  // After the ungz stream has ended, verify the checksum.
  ungz
    .on("end", () => {
      const expectedChecksum = checksumMap?.[pkgNameWithPlatform];
      // Skip verification if we can't find the file checksum
      if (!expectedChecksum) {
        console.warn("Skipping checksum verification");
        return;
      }
      const calculatedChecksum = hash.digest("hex");
      if (calculatedChecksum !== expectedChecksum) {
        throw errChecksum;
      }
      console.error("Checksum verified.");
    })
    // Pipe the data to the untar stream.
    .pipe(untar);

  // Wait for the untar stream to finish.
  await new Promise((resolve, reject) => {
    untar.on("error", reject);
    untar.on("end", () => resolve());
  });

  // Link the binaries in postinstall to support yarn
  await binLinks({
    path: path.resolve("."),
    pkg: { ...pkg, bin: { [Object.keys(pkg.bin)[0]]: binPath } },
  });

  console.error("Installed Studio MCP CLI successfully");
}

await main();
