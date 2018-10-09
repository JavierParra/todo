from pprint import pprint
import falcon


class Todos():
    def __init__(self):
        pass

    def on_get(self, req, resp):
        todo = {
            "id": 1,
            "complete": True,
            "name": "Take out the trash",
            "created": 1455151212,
            "completed": 1455171642
        }
        print(todo)
        resp.status = falcon.HTTP_200
        resp.body = 'Holas mundos'


todos = Todos()

app = falcon.API()

app.add_route('/todos', todos)
