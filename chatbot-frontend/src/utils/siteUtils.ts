/**
 * Convert a file path to a URL on the Epistemic Technology site
 * @param filePath - The file path to convert
 * @returns The URL
 */
export const filePathToURL = (filePath: string) => {
  // Extract the URL path from the file path by removing the base directory and changing .md to /
  const baseDir = "/Users/mikethicke/github/epistemic.technology/site/content/";
  let urlPath = filePath;
  if (filePath.startsWith(baseDir)) {
    const relativePath = filePath.substring(baseDir.length);
    urlPath = relativePath.replace(/\.md$/, "/");
  }
  return new URL(urlPath, window.location.origin).toString();
};
