# Epistemic Technology Website

This is the main website for Epistemic Technology, built using Hugo static site generator. It serves as the company's online presence and integrates with various backend services.

## Architecture

The site is built with [Hugo](https://gohugo.io/) and uses a custom theme (ETBASIC). It's designed to be served as a static site via NGINX, with several backend services providing dynamic functionality:

- Contact form submissions are handled by the [contact-backend](../contact-backend/README.md)
- Blog subscriptions are processed by the [blog-signup](../blog-signup/README.md) service
- Interactive chat functionality is provided by the [chatbot-frontend](../chatbot-frontend/README.md) and [chatbot-backend](../chatbot-backend/README.md)

The site is containerized using Docker.

## Build Parameters

The site build process accepts several build arguments to configure API endpoints:

- `HUGO_PARAMS_formSubmitEndpoint` - Endpoint for the contact form API
- `HUGO_PARAMS_subscribeSubmitEndpoint` - Endpoint for the blog subscription signup API
- `HUGO_PARAMS_subscribeConfirmEndpoint` - Endpoint for the blog subscription confirmation API
- `VITE_API_URL` - Endpoint for the chatbot API
- `VITE_FILEPATH_BASE_DIR` - Base directory path for resolving file paths in the chatbot

## Running

The site is expected to be run through Docker, which handles building both the Hugo site and the chatbot frontend, then combines them into a single NGINX-served container.

```bash
docker build -t epistemic-site \
  --build-arg HUGO_PARAMS_formSubmitEndpoint=https://api.epistemic.technology/contact/ \
  --build-arg HUGO_PARAMS_subscribeSubmitEndpoint=https://api.epistemic.technology/blog/signup \
  --build-arg HUGO_PARAMS_subscribeConfirmEndpoint=https://api.epistemic.technology/blog/confirm \
  --build-arg VITE_API_URL=https://api.epistemic.technology/chat \
  -f site/Dockerfile .

docker run -p 80:80 epistemic-site
```

For local development, you can run Hugo directly:

```bash
cd site
hugo server
```
