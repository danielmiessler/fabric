from app import create_app, db
from app.models import User, Pattern
import os
from app import sockio

app = create_app()
current_directory = os.path.dirname(os.path.realpath(__file__))


@app.cli.command('init-db')
def init_db_command():
    with app.app_context():
        db.create_all()
        # if User.query.filter_by(username='username').first() is None:
        #     user = User(username='username')
        #     user.set_password('password')
        #     user.is_admin = True
        #     db.session.add(user)
        #     db.session.commit()
        baseline_patterns = os.listdir(os.path.join(
            current_directory, 'app/chatgpt/patterns'))
        for pattern in baseline_patterns:
            if Pattern.query.filter_by(name=pattern).first() is None:
                with open(os.path.join(current_directory, f'app/chatgpt/patterns/{pattern}/system.md'), 'r') as f:
                    pattern_text = f.read()
                new_pattern = Pattern(name=pattern, pattern=pattern_text)
                db.session.add(new_pattern)
                db.session.commit()


if __name__ == '__main__':
    sockio.run(app)
