# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm install

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Run stage
FROM nginx:alpine

# Remove default nginx config
RUN rm /etc/nginx/conf.d/default.conf

# Copy nginx config template
COPY nginx.conf.template /etc/nginx/templates/default.conf.template

# Copy built assets from builder
COPY --from=builder /app/dist /usr/share/nginx/html

# Railway provides PORT env var — default to 80 for local
ENV PORT=80

# Nginx's docker-entrypoint uses envsubst on /etc/nginx/templates/*.template
# and outputs to /etc/nginx/conf.d/ automatically

EXPOSE ${PORT}

CMD ["nginx", "-g", "daemon off;"]
