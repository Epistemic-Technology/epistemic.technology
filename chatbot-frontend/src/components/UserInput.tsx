import { Component, createSignal, createEffect, batch } from "solid-js";
import { UserInputProps } from "../types";
import LoadingSpinner from "./LoadingSpinner";

const UserInput: Component<UserInputProps> = (props) => {
  const [inputValue, setInputValue] = createSignal("");
  const [formSubmitted, setFormSubmitted] = createSignal(false);
  const [isFocused, setIsFocused] = createSignal(false);
  let inputRef: HTMLInputElement | undefined;

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    const query = inputValue().trim();
    if (query && !props.isLoading) {
      setInputValue("");

      try {
        await props.onSubmit(query);
      } finally {
        setFormSubmitted(true);
      }
    }
  };

  createEffect(() => {
    if (formSubmitted() && inputRef && !props.isLoading) {
      inputRef.focus();
      setFormSubmitted(false);
    }
  });

  createEffect(() => {
    if (props.autoFocus && inputRef) {
      inputRef.focus();
    }
  });

  return (
    <form onSubmit={handleSubmit} class="flex items-center">
      <span>&gt;&nbsp;</span>
      <div class="relative flex-grow user-input-container">
        {props.isLoading ? (
          <div class="flex items-center">
            <LoadingSpinner />
            <span>Processing your request...</span>
          </div>
        ) : (
          <>
            <input
              type="text"
              value={inputValue()}
              onInput={(e) => setInputValue(e.currentTarget.value)}
              onFocus={() => setIsFocused(true)}
              onBlur={() => setIsFocused(false)}
              disabled={props.isLoading}
              class="w-full bg-transparent border-none focus:outline-none caret-transparent !text-2xl user-input"
              autocomplete="off"
              ref={inputRef}
            />
            <div
              class="cursor-block text-2xl"
              style={{
                left: `${inputValue().length}ch`,
                display: props.isLoading || !isFocused() ? "none" : "block",
              }}
            ></div>
          </>
        )}
      </div>
      <style>
        {`
          .cursor-block {
            position: absolute;
            top: 0;
            width: 0.6em;
            height: 1.2em;
            background-color: currentColor;
            animation: blink 1s step-end infinite;
          }
          
          @keyframes blink {
            0%, 100% { opacity: 1; }
            50% { opacity: 0; }
          }
        `}
      </style>
    </form>
  );
};

export default UserInput;
