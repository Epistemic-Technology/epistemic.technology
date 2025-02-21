---
title: "Building My Site with Cursor"
date: 2025-02-23
draft: false
---

[Cursor](https://www.cursor.com/) is a code editor based on Microsoft's Visual Studio Code. Initially released in 2023, it has become one of the leading AI-focused code editors. I've been using Cursor for a while now but hadn't tried out some of its newer features like Agent Mode. So building a new site from scratch seemed like the perfect opportunity to experiment.

## Initial Planning and Architecture

For any new project, coding or otherwise, my starting point is ChatGPT. You could use Cursor's chat mode for this, but I still find OpenAI's chat interface better for developing my initial thoughts. Here's my initial prompt:

![Asking ChatGPT: "I'm making a website. It's mostly static content. A home page, an about page, a contact page with a form, and a blog section. The blog content could be static html compiled from markdown, or something similar.  The one piece of dynamic content is a chat interface that could be a React app or vanilla js or whatever. When the user chats, the query will be posted to a backend endpoint that does RAG. So a vector database that could be in sqllite or something else lightweight. The retrieved document will be added to context for an API call to GPT or Claude, etc. I'd prefer to write the backend code in Go and the blog content in markdown. Can you suggest an approach and stack?"](/images/blog/2025-02-23-building-my-site-with-cursor/chatgpt-stack.png)

After some back-and-forth, we decided on the following stack:

- [Hugo](https://gohugo.io/) for generating the static site content.
- Go for backend services such as processing the contact form and blog subscriptions.
- Vanilla JavaScript for functionality I couldn't achieve with HTML & CSS alone.

The chatbot functionality on the site didn't make it into the initial launch, but it will be Go, SQLite, and Svelte.

Everything is deployed using [Railway](https://railway.com/), which is such a breath of fresh air after spending much of the last four years of my life struggling with AWS.

My philosophy toward web development has evolved significantly over the last few years from believing that you should just do everything in a full-featured framework like NextJS or Laravel to my current stance of trying to find the absolute simplest, most lightweight solution I can find for the current problem. For Epistemic Technology, which doesn't need to be richly interactive, having a mostly static site eliminates so many problems. I don't have to worry about server-side vs. client-side rendering. I'm using vanilla CSS and JS, so there is no build pipeline beyond Hugo compiling the site from Markdown and template files. There's no database! The site loads quickly and (I think) feels like a single page app even though it isn't. And it does really well on Google's PageSpeed Insights:

![Google PageSpeed insights showing 100 Performance, 100 Accessibility, 100 Best Practices, and 91 SEO scores](/images/blog/2025-02-23-building-my-site-with-cursor/pagespeed.png)

Yes, it took some work to get this point, but keeping things simple made it way easier.

## Cursor Interface and Setup

![The Cursor interface, with the main editor on the left, Composer in the right sidebar, and a terminal window below](/images/blog/2025-02-23-building-my-site-with-cursor/cursor.png)

Cursor should be familiar to anyone who has used VS Code or pretty much any modern IDE. The main differences between Cursor and VS Code are that Cursor is more focused on an AI-first workflow, where you are constantly working in the right-hand sidebar, either chatting with a model to think through ideas or gain understanding, or giving it directions in Composer to make changes or additions to your project. Cursor's UI can be overwhelming and is rapidly evolving, but it's really powerful. The only real customization I made was to make ⌘N open a new Composer session so I can more easily jump back and forth between the editor window and the sidebar without having to take my hands off the keyboard.

Cursor lets you choose which model to use for queries. Like most developers these days, I'm using claude-3.5-sonnet, though in my experience there isn't much difference between any of the frontier models.

## Building the static site

I'm new to Hugo, so a big hurdle to get over at the beginning was just figuring out how to do *anything* in a new framework. Where are images supposed to go? How do templates work? How can pass parameters to a component? Here the Chat interface was the most useful. For example:

![Chat sidebar interaction. I ask "In hugo, where do I put pages, like an about page?" Cursor responds by explaining Hugo's directory structure and giving an example of metadata for a specific page or post.](/images/blog/2025-02-23-building-my-site-with-cursor/cursor-chat.png)

One of the confusing parts of Cursor is that pretty much any interaction can result in code generation. So if it makes a suggestion in chat you can apply that suggestion to your code. In Composer, you can use "normal" mode or "agent" mode. And you can interact with the model from the editor window as well to have it make changes or to ask about a specific portion of the file. There are subtle differences between each of these modes. Agent mode is the most active in making project-wide changes, interacting with the terminal, downloading new packages, and so on. "Normal" Composer mode is more restricted in both what it can do and the context in which it operates (you need to point it at specific files), and chat is more restrictive still. That can be overwhelming and I feel like Cursor is going to be doing a lot of iteration on its user experience.

One of the first bumps in the road I encountered with Cursor was generating the CSS for the site. I'm using vanilla CSS, mostly in a single file. I'm not primarily a frontend developer, and CSS has a million little intricacies, so I had to lean on the AI pretty heavily to style the site. I gave it a starting point by setting up some variables for dark and light mode colors and picking fonts. Then I mostly gave it general directions about how to style the site. However, things got gnarly quickly as, for instance, all sorts of random elements had `2em` left and right padding to create space with the edge of the window rather than putting that padding on the `main` section.

I think two things are going on here. First, LLMs still don't have human-like understanding of code. When you tell them to do something, they're pretty good at satisfying the request, but they aren't going to step back like a human would and say, "You know, I think it would be better if...". Second, things can go off the rails if you make too many incremental changes in a row. Again, the model doesn't have the holistic view and understanding of a human, and so it's not going to pause in fulfilling your request to do some much-needed refactoring.

## Building the backend services

For the initial site launch, I wanted two backend services to extend its functionality: a contact form and a blog subscription form. For both of these, I could have relied entirely on external services and embedded forms in `<iframe>`s, but I wanted to have more control and a better user experience. So I made a couple of Go services that would accept form submissions and take appropriate action. Ultimately, I'm using an external service for delivering email (Sendgrid), but these backend services mediate between the site and that external service and handle things like rate limiting, confirming email addresses, and so on.

![Screenshot of email subscription confirmation, showing that the subscription was confirmed](/images/blog/2025-02-23-building-my-site-with-cursor/email-confirmation.png)

Here, I relied heavily on Agent mode, and it went surprisingly well. These are really simple services, a couple hundred lines of code each. I got mostly working versions of each on the first try and didn't make any major changes. In places I had to change how environment variables worked, or specific details about the requests and responses, but overall it was a really smooth process. I think having such a small scope helped a lot here, and it also might help to be writing in Go, which is a famously straightforward language. One of the design goals of Go was to have everyone's code look the same, from a junior developer to someone who has years of experience. This is a perfect scenario for an LLM, as it means that there will be a close correspondence between its training data and what it needs to produce.

It does make some errors with specific libraries that might not be as popular, like Sendgrid's Go library, but overall it was a great experience. It was *much less* successful with my next project, building a RAG-based chatbot for the site, which is significantly more complex. But that is for another post.

## What I learned

I've been using Cursor for a while, but this was my first project where I leaned heavily on its more agentic features. Here are a few things I've taken away:

- Cursor's Agent mode is way better than I expected. It's fast and strikes a good balance between autonomous action and keeping you in the loop.
- LLMs work best working in a narrow scope. "Agentic" AI isn't at the point where you can let it lose on a large codebase and have it succeed autonomously.
- LLMs add a workflow dependency, in that you rely on this thing that you never had to before, but they can take away other dependencies. For this initial stage I used vanilla JavaScript and CSS partly because I didn't need to use something like Tailwind or SASS to improve the developer experience. I could just have the model remember syntax.
- LLMs are great at big-picture tasks, like brainstorming your architecture, and small-scale tasks, like implementing a well-defined function. There is a huge middle zone where humans are still essential for understanding how to structure an application.
- The more specific you are in your prompts, the better result you will get out of the model. I had a very well-defined idea of what I wanted out of the site, what technologies I wanted to use, and how different components would connect. I don't think you can realistically tell an "agent" to build an app from some a general, vague description and hope to get something useful out of it. We are nowhere near a world where software engineers will be obsolete.

One question that I am still bouncing around in my head is how tools like Cursor should affect how we design our systems. I think the expectation of most people right now is that agentic AI will design software like we do, but better, or faster, or cheaper. But I think what we should be doing is architecting with AI in mind from the beginning: doing system design in a way that facilitates AI-assisted development.

But that's for another post.

I hope you got something out of this. If you did, please consider sharing it and subscribing!