# Use an official Node.js runtime as a parent image
FROM node:20

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application
COPY . .

# Build the Next.js application
#RUN npm run build

# Expose port 3000 (Next.js default)
EXPOSE 3000

# Start the Next.js application
CMD ["npm", "run", "dev"]
