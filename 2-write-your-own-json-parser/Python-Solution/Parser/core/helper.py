
class Helper:
    def get_bytes(file_name):
        file_stat = os.stat(file_name)
        return file_stat.st_size

    def get_lines(file_name=None):
        return sum(1 for _ in open(file_name))

    def get_characters(file_name=None):
        return sum(len(word) for word in open(file_name))

    def get_words(file_path=None):
        file = open(file_path, 'r', encoding='utf-8')
        return len(file.read().split())

    def tokens():
        identifier = []
        keyword = ["and", "del", "from", "not", "while", "as", "elif", "global", "or", "with", "assert", "else", "if", "pass", "yield", "break",
                   "except", "import", "print", "class", "exec", "in", "raise", "continue", "finally", "is", "return", "def", "for", "lambda", "try"]
        separator = ["\newline", "\\", "\'", '\"', "\a", "\b", "\f", "\n", '\N"{name}', "\r", "\t", "\uxxxx", "\Uxxxxxxxx", "\v", "\ooo", "\xhh"]
        operator = ["+", "-", "*", "**", "//",  "<<", ">",
                    "|", "^", "~", "<", ">", "<=", "=", "=", "=", ">"]
        literals = []
        comment = []
