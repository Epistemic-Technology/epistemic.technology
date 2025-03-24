import { Component, For, JSX } from "solid-js";
import { BotMessageProps } from "../types";
import TerminalMessageWrapper from "./TerminalMessageWrapper";
import { filePathToURL } from "../utils/siteUtils";

const BotMessage: Component<BotMessageProps> = (props) => {
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
                    <span class="mr-1 pr-2">
                      {source.index ? `:${source.index}` : ""}
                    </span>
                    <a
                      class="!text-blue-300"
                      href={filePathToURL(source.FilePath)}
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
