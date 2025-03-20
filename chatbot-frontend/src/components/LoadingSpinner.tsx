import { Component, createSignal, onCleanup, onMount } from "solid-js";

const LoadingSpinner: Component = () => {
  const frames = ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"];
  const [frame, setFrame] = createSignal(0);
  let intervalId: number;

  onMount(() => {
    intervalId = window.setInterval(() => {
      setFrame((prev) => (prev + 1) % frames.length);
    }, 100);
  });

  onCleanup(() => {
    clearInterval(intervalId);
  });

  return (
    <span
      class="text-3xl inline-block animate-pulse"
      aria-live="polite"
      aria-label="Loading"
    >
      {frames[frame()]}
    </span>
  );
};

export default LoadingSpinner;
