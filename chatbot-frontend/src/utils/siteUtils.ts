/**
 * Convert a file path to a URL on the Epistemic Technology site
 * @param filePath - The file path to convert
 * @returns The URL
 */
export const filePathToURL = (filePath: string) => {
  // Extract the URL path from the file path by removing the base directory and changing .md to /
  const baseDir = import.meta.env.VITE_FILEPATH_BASE_DIR;
  let urlPath = filePath;
  console.log("baseDir: ", baseDir);
  console.log("filePath: ", filePath);
  if (filePath.startsWith(baseDir)) {
    const relativePath = filePath.substring(baseDir.length);
    urlPath = relativePath.replace(/\.md$/, "/");
  }
  return new URL(urlPath, window.location.origin).toString();
};
