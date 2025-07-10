
---
title: "Projects"
description: "Projects completed by Epistemic Technology"
keywords: "AI development, educational technology, research technology, web development"
date: 2025-07-09
draft: false
---

## Epistemic Technology Projects

*Projects I've worked on within Epistemic Technology*

### Epistemic Technology Website

![Screenshot of the Epistemic Technology homepage showing hero image and mission statement](/images/projects/epistemic-technology-home.png)

- Building and deploying the Epistemic Technology website.
- **Dates:** January - February 2025
- **Technologies:** Hugo, Go, HTML, CSS, JavaScript, Railway, SendGrid
- **GitHub:** [Epistemic-Technology/epistemic.technology](https://github.com/Epistemic-Technology/epistemic.technology)
- The Epistemic Technology website is a static website built using Hugo, a popular static site generator written in Go. It is hosted on Railway, a cloud platform for deploying and scaling web applications. The website uses SendGrid for email delivery. It has minimal JavaScript to handle form submissions and smooth scrolling. A few small backend services written in Go are used to handle functionality such as signing users up for blog subscriptions and handling contact form submissions.
- I discuss this project in my blog post [Building my site with Cursor](/blog/2025-02-23-building-my-site-with-cursor/).

### Epistemic Technology Chatbot

![Screenshot of the Epistemic Technology chatbot showing a brief interaction and the WarGames aesthetic](/images/projects/epistemic-technology-chatbot.png)

- Building a RAG-based chatbot for the Epistemic Technology website.
- **Dates:** March 2025
- **Technologies:** Go, OpenAI Go SDK, SQLite, SolidJS, TailwindCSS
- **GitHub:** [Epistemic-Technology/epistemic.technology](https://github.com/Epistemic-Technology/epistemic.technology)
- The Epistemic Technology chatbot is a RAG-based chatbot built using the OpenAI Go SDK for generating embeddings and responses, and SQLite as a vector database. It uses SolidJS for the frontend. The app is rendered on top of the static website, and TailwindCSS is used to isolate the chatbot's styles from the rest of the website. The chatbot's aesthetic is meant to evoke the classic 80s movie WarGames.
- I discuss this project in my blog post [Building a Chatbot](/blog/2025-03-31/building-a-chatbot/)

### Co-Intelligence AI Obsidian Plugin

![Screenshot of the Co-Intelligence AI Obsidian Plugin, showing a text-based interface similar to ChatGPT](/images/projects/obsidian-co-intelligence.png)

- Building a multi-provider chat application for [Obsidian](https://obsidian.md)
- **Dates:** April - May 2025
- **Technologies:** TypeScript, Vercel AI SDK, Obsidian API
- **GitHub:** [Epistemic-Technology/co-intelligence](https://github.com/Epistemic-Technology/co-intelligence)
- Co-Intelligence provides a full chat interface within Obsidian, allowing users to interact with a variety of LLMs, provide context from notes in their vault, and automatically save conversations as markdown-formatted notes.

### GitHub Tasks Obsidian Plugin

![Screenshot of the GitHub Tasks Obsidian Plugin showing an Obsidian note with a list of tasks](/images/projects/obsidian-github-tasks.png)

- Building an Obsidian plugin to sync GitHub issues and pull requests to an Obsidian note.
- **Dates:** June 2025
- **Technologies:** TypeScript, GitHub API, Obsidian API
- **GitHub:** [Epistemic-Technology/obsidian-github-tasks](https://github.com/Epistemic-Technology/obsidian-github-tasks)
- GitHub Tasks is a simple, but powerful and flexible tool for synchronizing GitHub issues and pull requests with Obsidian tasks. I made it mainly so that I could use it myself to see all of my outstanding tasks in a single place.

### Museum for WordPress Plugin

![Screenshot of the Museum for WordPress Plugin running on utsic.utoronto.ca](/images/projects/utsic.png)

- Building a full-featured museum management system in WordPress.
- **Dates:** 2017-2021; 2024-present (major overhaul)
- **Technologies:** PHP, WordPress API, JavaScript, SCSS, React
- **GitHub:** [Epistemic-Technology/wp-museum](https://github.com/Epistemic-Technology/wp-museum)
- Museum for WordPress is a full-featured museum management system built for the [University of Toronto Scientific Instruments Collection](https://www.utsic.utoronto.ca/). It is designed to make it easy for organizations and individuals to manage and share their collections of physical artifacts while tightly integrating with a WordPress website.

## Other Projects

*Projects I've worked on beyond Epistemic Technology*

### Knowledge Commons Search

![Screenshot of the Knowledge Commons Search architecture diagram](/images/projects/cc-search.png)

- Building a backend search service for the Knowledge Commons network.
- **Dates:** January - May 2024
- **Technologies:** Go, ElasticSearch / OpenSearch, PHP
- **GitHub:** [MESH-Research/commons-connect](https://github.com/MESH-Research/commons-connect)
- As the [Knowledge Commons](https://kcommons.org/) network expanded from a monolithic WordPress instance to a network of applications, we needed a new backend search that could unify users' content. I built a search service in Go that serves an API gateway backed by Amazon's OpenSearch service. Then I built WordPress functionality that provisioned users content to the service and queried the service for results.

### Knowledge Commons Containerization

![Screenshot of the Amazon ECS dashboard showing the Knowledge Commons WordPress services](/images/projects/wordpress-ecs.png)

- Modernizing and containerizing the Knowledge Commons WordPress application.
- **Dates:** May 2023 - July 2024
- **Technologies:** Docker, Lando, AWS Elastic Container Service (ECS), AWS Elastic Container Registry (ECR), AWS Secrets Manager, GitHub Actions
- **GitHub:** [MESH-Research/knowledge-commons-wordpress](https://github.com/MESH-Research/knowledge-commons-wordpress)
- Before this project, [Knowledge Commons](https://kcommons.org/) ran on a traditional LEMP stack, with development taking place through live filesystem edits and pulling changes from GitHub. Local development was infeasible. This (very large) project involved centralizing the Knowledge Commons WordPress application in a single repository, containerizing the application, using Lando for local development, moving most configuration to environment variables configured through AWS Secrets Manager, creating a CI/CD pipeline with GitHub Actions, and deploying to AWS ECS with a blue/green deployment strategy that allows for zero downtime deployments and end-to-end testing of deployments before exposing them to regular users.

### Knowledge Commons Mailchimp Integration

![Screenshot of an Humanities Commons onboarding email](/images/projects/mailchimp.png)

- Building Mailchimp integration for the Knowledge Commons network.
- **Dates:** September 2023
- **Technologies:** PHP, WordPress API, Mailchimp API
- **GitHub:** [MESH-Research/knowledge-commons-wordpress](https://github.com/MESH-Research/knowledge-commons-wordpress)
- Knowledge Commons is a complex social networking application, and one of its persistent challenges is onboarding new users. Previously we kept in touch with users through our transactional email system, but it was clunky and difficult for our community team to manage. I built the Mailchimp integration to automatically kick off a series of onboarding emails to new users and to add and remove users from our regular newsletter.

### Knowledge Commons HASTAC Migration

![Screenshot of the HASTAC homepage](/images/projects/hastac.png)

- Migrating a large Drupal-based academic social network to the Knowledge Commons WordPress platform.
- **Dates:** February - November 2022
- **Technologies:** PHP, WordPress API, Drupal API
- HASTAC (Humanities, Arts, Science, and Technology Alliance and Collaboratory) is a venerable academic social network that needed a new, more sustainable home. Working in collaboration with their developer, I migrated thousands of posts, pages, and sites, and approximately 20,000 user accounts to Knowledge Commons. The involved programmatically creating user accounts, sites, and posts, linking everything together, and creating a mechanism for users to retain control of their content on a new platform.
