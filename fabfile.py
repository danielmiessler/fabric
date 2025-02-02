from fabric2 import task
import time

@task
def start_backend(ctx):
    """Start the backend server"""
    ctx.run("python backend/app.py", asynchronous=True)

@task
def start_frontend(ctx):
    """Start the Streamlit frontend"""
    ctx.run("streamlit run frontend/streamlit_app.py", asynchronous=True)

@task
def start_all(ctx):
    """Start both backend and frontend servers"""
    start_backend(ctx)
    # Wait for backend to initialize
    time.sleep(3)
    start_frontend(ctx)

@task
def check_env(ctx):
    """Check environment and dependencies"""
    ctx.run("python --version")
    ctx.run("pip list")
    ctx.run("netstat -tulpn | grep LISTEN")
