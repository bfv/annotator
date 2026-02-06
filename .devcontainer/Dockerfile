FROM mcr.microsoft.com/devcontainers/go:dev-1.25-bookworm

# Verwijder de Yarn repository die de GPG-fout veroorzaakt. 
# De Docker-feature installatie zal daarna zonder fouten apt-get update kunnen draaien.
RUN rm -f /etc/apt/sources.list.d/yarn.list \
    && rm -f /etc/apt/sources.list.d/nodesource.list \
    && apt-get update
    