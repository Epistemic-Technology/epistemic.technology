import { Component, createSignal, createEffect, batch } from "solid-js";
import { UserInputProps } from "../types";

const UserInput: Component<UserInputProps> = (props) => {
  const [inputValue, setInputValue] = createSignal("");
  const [formSubmitted, setFormSubmitted] = createSignal(false);
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

  // Focus management
  createEffect(() => {
    if (formSubmitted() && inputRef && !props.isLoading) {
      inputRef.focus();
      setFormSubmitted(false);
    }
  });

  // Initial focus
  createEffect(() => {
    if (props.autoFocus && inputRef) {
      inputRef.focus();
    }
  });

  return (
    <form onSubmit={handleSubmit} class="flex items-center">
      <span>&gt;&nbsp;</span>
      <div class="relative flex-grow">
        <input
          type="text"
          value={inputValue()}
          onInput={(e) => setInputValue(e.currentTarget.value)}
          disabled={props.isLoading}
          class="w-full bg-transparent border-none focus:outline-none caret-transparent text-3xl"
          autocomplete="off"
          ref={inputRef}
        />
        <div
          class="cursor-block"
          style={{
            left: `${inputValue().length}ch`,
            display: props.isLoading ? "none" : "block",
          }}
        ></div>
      </div>
      <style>
        {`
          .cursor-block {
            position: absolute;
            top: 0;
            width: 0.8em;
            height: 1.5em;
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
