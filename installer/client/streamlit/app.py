import streamlit as st
from utils import Standalone, Update, Alias, Setup, Transcribe
import sys, os, io
from contextlib import redirect_stdout


class Args:
    def __init__(self):
        self.model = None
        self.temp = 0.7
        self.top_p = 0.9
        self.frequency_penalty = 0
        self.presence_penalty = 0
        self.copy = False
        self.output = ""
        self.session = ""
        self.listmodels = False


def fetch_available_models(standalone):
    gptmodels, localmodels, claudemodels, googlemodels = (
        standalone.fetch_available_models()
    )
    return gptmodels, localmodels, claudemodels, googlemodels


def main():
    st.set_page_config(
        page_title="Fabric",
        page_icon=":robot_face:",
        initial_sidebar_state="expanded",
        layout="centered",
    )

    st.title("Fabric")
    st.markdown("---")

    current_directory = os.path.dirname(os.path.realpath(__file__))
    config_directory = os.path.expanduser("~/.config/fabric")
    env_file = os.path.join(config_directory, ".env")

    home_holder = os.path.expanduser("~")
    config = os.path.join(home_holder, ".config", "fabric")
    config_patterns_directory = os.path.join(config, "patterns")
    config_context = os.path.join(config, "context.md")
    env_file = os.path.join(config, ".env")

    # Initialize Standalone instance
    args = Args()
    standalone = Standalone(args)

    gptmodels, localmodels, claudemodels, googlemodels = fetch_available_models(
        standalone
    )

    models = ()

    if gptmodels:
        for model in gptmodels:
            models += (model,)

    if localmodels:
        for model in localmodels:
            models += (model,)

    if claudemodels:
        for model in claudemodels:
            models += (model,)

    if googlemodels:
        for model in googlemodels:
            models += (model,)

    with st.sidebar:
        st.title("Settings")
        selected_model = st.selectbox(
            "Models", options=models, index=None, placeholder="Select a model"
        )

        try:
            direct = sorted(os.listdir(config_patterns_directory))
            if direct:
                selected_pattern = st.selectbox(
                    "Patterns",
                    options=direct,
                    index=None,
                    placeholder="Select a pattern",
                )
            else:
                st.warning("No patterns found")
        except FileNotFoundError:
            st.warning("No patterns found")

        temp = st.slider(
            "Temperature", min_value=0.0, max_value=1.0, value=0.7, step=0.1
        )
        top_p = st.slider("Top P", min_value=0.0, max_value=1.0, value=0.9, step=0.1)

    user_input = st.chat_input(
        "Paste YouTube link or ask something", max_chars=None, key="user_input"
    )

    transcribe = Transcribe()

    if user_input:
        args.model = selected_model
        args.temp = temp
        args.top_p = top_p

        with st.chat_message("user"):
            st.markdown(user_input)

            # Initialize Standalone instance with Args object
        args = Args()
        standalone = Standalone(args=args, pattern=selected_pattern, env_file=env_file)

        output = io.StringIO()
        with st.chat_message("assistant"):
            with st.spinner("Generating response..."):
                with redirect_stdout(output):
                    standalone.streamMessage(user_input, context="", host="")
            response = output.getvalue()

            # Display the response
            st.write(response)

    update = Update()
    alias = Alias()
    with st.sidebar:
        if st.button("Update Patterns", use_container_width=True):
            update.update_patterns()
            alias.execute()
            st.success("Patterns updated successfully!")

        if st.button(
            "Default Model",
            help="Set the selected model from above as default model",
            use_container_width=True,
        ):
            setup = Setup()
            try:
                setup.default_model(selected_model)
                st.success(f"Default model changed to {selected_model} successfully!")
                st.info("Please restart the app for changes to take effect!")
            except Exception as e:
                st.error(f"Error: {e}")


if __name__ == "__main__":
    main()
