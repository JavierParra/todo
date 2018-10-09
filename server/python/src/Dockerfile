FROM python

RUN mkdir -p /usr/src
WORKDIR /usr/src

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY . .

EXPOSE 8000

CMD ["gunicorn", "--bind=0.0.0.0:8000", "--reload", "app.app:app"]
