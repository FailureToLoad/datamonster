FROM node:22.2-alpine 
ENV NODE_ENV development
WORKDIR /app

COPY ./ ./

RUN npm ci --no-audit
CMD [ "npm", "run", "dev" ]