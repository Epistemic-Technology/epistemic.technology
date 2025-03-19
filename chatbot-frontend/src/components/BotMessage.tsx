import { Component, For, JSX } from "solid-js";
import { BotMessageProps } from "../types";
import TerminalMessageWrapper from "./TerminalMessageWrapper";

const BotMessage: Component<BotMessageProps> = (props) => {
  const filePathToURL = (filePath: string) => {
    // Extract the URL path from the file path by removing the base directory and changing .md to /
    const baseDir =
      "/Users/mikethicke/github/epistemic.technology/site/content/";
    let urlPath = filePath;
    if (filePath.startsWith(baseDir)) {
      const relativePath = filePath.substring(baseDir.length);
      urlPath = relativePath.replace(/\.md$/, "/");
    }
    return new URL(urlPath, window.location.origin).toString();
  };

  return (
    <TerminalMessageWrapper>
      <div>
        {props.content.split("\n").map((paragraph, index) => (
          <p class="mt-2 !p-0">{paragraph || "\u00A0"}</p>
        ))}
        {props.sources && (
          <>
            <p>Sources:</p>
            <ol class="list-none">
              <For each={props.sources}>
                {(source, index) => (
                  <li class="ml-2">
                    <span class="mr-1 pr-2">:{index() + 1}</span>
                    <a
                      class="!text-blue-300"
                      href={filePathToURL(source.FilePath)}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      {source.Title}
                    </a>
                  </li>
                )}
              </For>
            </ol>
          </>
        )}
      </div>
    </TerminalMessageWrapper>
  );
};

export default BotMessage;
