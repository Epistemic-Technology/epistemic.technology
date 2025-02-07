FROM hugomods/hugo:base AS builder
COPY . .
RUN hugo build
FROM nginx:alpine AS runner
COPY --from=builder /src/public /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
COPY default.conf.template /etc/nginx/templates/default.conf.template
