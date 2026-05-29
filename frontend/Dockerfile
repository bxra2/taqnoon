# Stage 1: Build the SvelteKit application
FROM node:22-alpine AS builder

WORKDIR /app

# Install dependencies (do this before copying source files for better caching)
COPY package.json package-lock.json* ./
RUN npm install

# Copy the rest of the application
COPY . .

# Build the application
RUN npm run build

# Prune non-production dependencies to make the image smaller
RUN npm prune --production

# Stage 2: Create the production image
FROM node:22-alpine

WORKDIR /app

# Copy the built application and production dependencies from the builder stage
COPY --from=builder /app/build ./build
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./

# Expose the port the app runs on
EXPOSE 3000

# Set environment variables for production
ENV NODE_ENV=production
ENV PORT=3000
# IMPORTANT: Provide ORIGIN environment variable in the Render Dashboard

# Start the Node.js server
CMD ["node", "build"]
