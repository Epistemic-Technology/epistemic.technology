FROM hugomods/hugo:base AS builder
COPY . .
ARG HUGO_PARAMS_formSubmitEndpoint=http://localhost:8080/
RUN echo $HUGO_PARAMS_formSubmitEndpoint
RUN hugo build
FROM nginx:alpine AS runner
COPY --from=builder /src/public /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
COPY default.conf.template /etc/nginx/templates/default.conf.template
RUN mkdir -p /data/logs
