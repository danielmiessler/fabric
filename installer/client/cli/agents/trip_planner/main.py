from crewai import Crew
from textwrap import dedent
from .trip_agents import TripAgents
from .trip_tasks import TripTasks
import os
from dotenv import load_dotenv
import agentops

# Load environment variables
current_directory = os.path.dirname(os.path.realpath(__file__))
config_directory = os.path.expanduser("~/.config/fabric")
env_file = os.path.join(config_directory, ".env")
load_dotenv(env_file)
os.environ['OPENAI_MODEL_NAME'] = 'gpt-4-0125-preview'

# Initialize AgentOps
AGENTOPS_API_KEY = os.getenv("AGENTOPS_API_KEY")
if not AGENTOPS_API_KEY:
    raise ValueError("AGENTOPS_API_KEY not found in environment variables")

agentops.init(AGENTOPS_API_KEY)

@agentops.track_agent(name='TripCrew')
class TripCrew:

    def __init__(self, origin, cities, date_range, interests):
        self.cities = cities
        self.origin = origin
        self.interests = interests
        self.date_range = date_range

    @agentops.record_function('run')
    def run(self):
        agents = TripAgents()
        tasks = TripTasks()

        city_selector_agent = agents.city_selection_agent()
        local_expert_agent = agents.local_expert()
        travel_concierge_agent = agents.travel_concierge()

        identify_task = tasks.identify_task(
            city_selector_agent,
            self.origin,
            self.cities,
            self.interests,
            self.date_range
        )
        gather_task = tasks.gather_task(
            local_expert_agent,
            self.origin,
            self.interests,
            self.date_range
        )
        plan_task = tasks.plan_task(
            travel_concierge_agent,
            self.origin,
            self.interests,
            self.date_range
        )

        crew = Crew(
            agents=[
                city_selector_agent, local_expert_agent, travel_concierge_agent
            ],
            tasks=[identify_task, gather_task, plan_task],
            verbose=True
        )

        result = crew.kickoff()
        return result

@agentops.track_agent(name='planner_cli')
class planner_cli:
    @agentops.record_function('ask')
    def ask(self):
        print("## Welcome to Trip Planner Crew")
        print('-------------------------------')
        location = input(
            dedent("""
        From where will you be traveling from?
        """))
        cities = input(
            dedent("""
        What are the cities options you are interested in visiting?
        """))
        date_range = input(
            dedent("""
        What is the date range you are interested in traveling?
        """))
        interests = input(
            dedent("""
        What are some of your high level interests and hobbies?
        """))

        trip_crew = TripCrew(location, cities, date_range, interests)
        result = trip_crew.run()
        print("\n\n########################")
        print("## Here is you Trip Plan")
        print("########################\n")
        print(result)

@agentops.record_function('main')
def main():
    planner = planner_cli()
    planner.ask()

if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        agentops.log_error(str(e))
    finally:
        agentops.end_session('Success')