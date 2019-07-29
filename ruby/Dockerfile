FROM ruby:2.6

# throw errors if Gemfile has been modified since Gemfile.lock
RUN bundle config --global frozen 1

WORKDIR /app

COPY Gemfile Gemfile.lock ./

RUN bundle install

COPY lib ./lib/
COPY .config.prod.ru ./

COPY html ./html/

WORKDIR /app

CMD ["bundle", "exec", "puma", ".config.prod.ru", "-p", "80", "-v", "--debug"]

EXPOSE 80
