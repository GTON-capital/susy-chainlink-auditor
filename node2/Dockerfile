FROM smartcontract/chainlink:1.2.0

# Create the chainlink node root path
RUN mkdir /chainlink

# Arg for api user email, with default value
ARG API_USER_EMAIL="user@example.com"
ENV API_USER_EMAIL=$API_USER_EMAIL

# Arg for api user password, with default value
ARG API_USER_PASSWORD="PA@SSword1234!567"
ENV API_USER_PASSWORD=$API_USER_PASSWORD

# Arg for node wallet password, with default value
ARG WALLET_PASSWORD="PA@SSword1234!567"
ENV WALLET_PASSWORD=$WALLET_PASSWORD

# Create chainlink node required values to initialize with
RUN echo $API_USER_EMAIL > /chainlink/.api
RUN echo $API_USER_PASSWORD >> /chainlink/.api
RUN echo $WALLET_PASSWORD > /chainlink/.password



ENV DATABASE_URL=postgresql://postgres:password@postgres:5432/chainlink2?sslmode=disable


ENV FEATURE_OFFCHAIN_REPORTING=true
ENV FEATURE_OFFCHAIN_REPORTING2=true
ENV P2P_NETWORKING_STACK=V1V2
ENV P2PV2_LISTEN_ADDRESSES=0.0.0.0:9901




ENV P2P_LISTEN_PORT=9902