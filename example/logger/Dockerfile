FROM ruby:2.6.1

COPY Gemfile Gemfile.lock ./
RUN bundle install -j4

COPY post_to_fluentd.rb ./

CMD ["bundle", "exec", "ruby", "post_to_fluentd.rb"]
