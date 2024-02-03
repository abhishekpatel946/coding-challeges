import json


class Lexer:
    def analyse(json_object: dict) -> list[str]:
        try:
            print(json_object)
            strinigy_json = json.dumps(json_object)

            for i in range(len(strinigy_json)):
                print(strinigy_json[i])
                # if

        except Exception as e:
            print(e)
