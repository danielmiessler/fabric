import os
import sys


class Session:
    def __init__(self):
        home_folder = os.path.expanduser("~")
        config_folder = os.path.join(home_folder, ".config", "fabric")
        self.sessions_folder = os.path.join(config_folder, "sessions")
        if not os.path.exists(self.sessions_folder):
            os.makedirs(self.sessions_folder)

    def find_most_recent_file(self):
        # Ensure the directory exists
        directory = self.sessions_folder
        if not os.path.exists(directory):
            print("Directory does not exist:", directory)
            return None

        # List all files in the directory
        full_path_files = [os.path.join(directory, file) for file in os.listdir(
            directory) if os.path.isfile(os.path.join(directory, file))]

        # If no files are found, return None
        if not full_path_files:
            return None

        # Find the file with the most recent modification time
        most_recent_file = max(full_path_files, key=os.path.getmtime)

        return most_recent_file

    def save_to_session(self, system, user, response, fileName):
        file = os.path.join(self.sessions_folder, fileName)
        with open(file, "a+") as f:
            f.write(f"{system}\n")
            f.write(f"{user}\n")
            f.write(f"{response}\n")

    def read_from_session(self, filename):
        file = os.path.join(self.sessions_folder, filename)
        if not os.path.exists(file):
            return None
        with open(file, "r") as f:
            return f.read()

    def clear_session(self, session):
        if session == "all":
            for file in os.listdir(self.sessions_folder):
                os.remove(os.path.join(self.sessions_folder, file))
        else:
            os.remove(os.path.join(self.sessions_folder, session))

    def session_log(self, session):
        file = os.path.join(self.sessions_folder, session)
        if not os.path.exists(file):
            return None
        with open(file, "r") as f:
            return f.read()

    def list_sessions(self):
        sessionlist = os.listdir(self.sessions_folder)
        find_most_recent_file_result = self.find_most_recent_file()
        if find_most_recent_file_result is not None:
            most_recent = find_most_recent_file_result.split("/")[-1]
            for session in sessionlist:
                with open(os.path.join(self.sessions_folder, session), "r") as f:
                    firstline = f.readline().strip()
                    secondline = f.readline().strip()
                    if session == most_recent:
                        print(f"{session} **default** \"{firstline}\n{secondline}\n\"")
                    else:
                        print(f"{session} \"{firstline}\n{secondline}\n\"")
        else:
            print('No files present in sessions directory')
