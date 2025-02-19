FROM hugomods/hugo:base AS builder
COPY . .
ARG HUGO_PARAMS_formSubmitEndpoint=http://localhost:8080/
ARG HUGO_PARAMS_subscribeConfirmEndpoint=http://localhost:8081/confirm
ARG HUGO_PARAMS_subscribeSubmitEndpoint=http://localhost:8081/signup
RUN echo $HUGO_PARAMS_formSubmitEndpoint
RUN echo $HUGO_PARAMS_subscribeConfirmEndpoint
RUN echo $HUGO_PARAMS_subscribeSubmitEndpoint
RUN hugo build
FROM nginx:alpine AS runner
COPY --from=builder /src/public /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
COPY default.conf.template /etc/nginx/templates/default.conf.template
RUN mkdir -p /data/logs
