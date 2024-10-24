FROM python:3.11-slim-buster

ENV APP_DIR=/
COPY requirements.txt $APP_DIR/requirements.txt
COPY run.sh $APP_DIR/run.sh
WORKDIR $APP_DIR

RUN python -m pip install --upgrade pip
RUN python -m pip install -r requirements.txt

CMD [ "/run.sh" ]
