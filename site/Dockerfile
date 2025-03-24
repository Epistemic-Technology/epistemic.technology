FROM hugomods/hugo:base AS builder
WORKDIR /app

# Install Node.js and npm
RUN apk add --no-cache nodejs npm

# Copy site content
COPY site/ ./site/
WORKDIR /app/site

# Build Hugo site
ARG HUGO_PARAMS_formSubmitEndpoint=http://localhost:8080/
ARG HUGO_PARAMS_subscribeConfirmEndpoint=http://localhost:8081/confirm
ARG HUGO_PARAMS_subscribeSubmitEndpoint=http://localhost:8081/signup
RUN hugo build

# Build chatbot frontend
WORKDIR /app/chatbot-frontend
COPY chatbot-frontend/ .
ARG VITE_API_URL=http://localhost:8082/chat
ARG VITE_FILEPATH_BASE_DIR=/app/site/content/
ENV VITE_API_URL=${VITE_API_URL}
ENV VITE_FILEPATH_BASE_DIR=${VITE_FILEPATH_BASE_DIR}
RUN npm install && npm run build

# Final stage
FROM nginx:alpine
COPY --from=builder /app/site/public /usr/share/nginx/html
COPY site/nginx.conf /etc/nginx/nginx.conf
COPY site/default.conf.template /etc/nginx/templates/default.conf.template
RUN mkdir -p /data/logs
COPY --from=builder /app/chatbot-frontend/dist /usr/share/nginx/html/chatbot

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
