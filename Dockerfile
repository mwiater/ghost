# Stage 1:
# Use base Alpine image to prepare our binary, label it 'app'
FROM golang:alpine AS app
# Install bash in this build stage (optional, only if needed during build)
RUN apk add --no-cache bash
# Add ghost user and group so that the Docker process in Scratch doesn't run as root
RUN addgroup -S ghost \
 && adduser -S -u 10000 -g ghost ghost
# Change to the correct directory to hold our application source code
WORKDIR /go/src/app
# Copy all the files from the base of our repository to the current directory defined above
COPY . .
# Compile the application to a single statically-linked binary file
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"' -tags timetzdata

# Stage 2:
# Use an Alpine base image instead of Scratch to include bash
FROM alpine:latest
# Install bash in the final stage
RUN apk add --no-cache bash
# Grab necessary certificates as Alpine has them, but update to ensure latest
RUN apk add --no-cache ca-certificates
# Copy our binary from the build stage to the Alpine image
COPY --from=app /go/bin/ghost /ghost
# Copy the user that we created in the first stage so that we don't run the process as root
COPY --from=app /etc/passwd /etc/passwd
# Change to the non-root user
USER ghost
# Run bash by default (you can change this if you'd like to keep the original entrypoint)
CMD ["/ghost"]
