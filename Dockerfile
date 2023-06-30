FROM --platform=linux/amd64 rustscan/rustscan:2.1.1 as rustscan
FROM --platform=linux/amd64 golang:1.19.5 as build
WORKDIR /temp
COPY . /temp/

RUN CGO_ENABLED=0 go build -v

FROM --platform=linux/amd64 alpine:3.12
COPY ./resource /resource
COPY --from=rustscan /usr/local/bin/rustscan /usr/local/bin/rustscan
COPY --from=build /temp/endgame /endgame
CMD [ "/endgame" ]