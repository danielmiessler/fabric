from crewai import Agent
from .tools.browser_tools import BrowserTools
from .tools.calculator_tools import CalculatorTools
from .tools.search_tools import SearchTools
import agentops
import os

# Initialize AgentOps
AGENTOPS_API_KEY = os.getenv("AGENTOPS_API_KEY")
if not AGENTOPS_API_KEY:
    raise ValueError("AGENTOPS_API_KEY not found in environment variables")

agentops.init(AGENTOPS_API_KEY)

@agentops.track_agent(name='TripAgents')
class TripAgents:

    @agentops.record_function('city_selection_agent')
    def city_selection_agent(self):
        return Agent(
            role='City Selection Expert',
            goal='Select the best city based on weather, season, and prices',
            backstory='An expert in analyzing travel data to pick ideal destinations',
            tools=[
                SearchTools.search_internet,
                BrowserTools.scrape_and_summarize_website,
            ],
            verbose=True)

    @agentops.record_function('local_expert')
    def local_expert(self):
        return Agent(
            role='Local Expert at this city',
            goal='Provide the BEST insights about the selected city',
            backstory="""A knowledgeable local guide with extensive information
        about the city, it's attractions and customs""",
            tools=[
                SearchTools.search_internet,
                BrowserTools.scrape_and_summarize_website,
            ],
            verbose=True)

    @agentops.record_function('travel_concierge')
    def travel_concierge(self):
        return Agent(
            role='Amazing Travel Concierge',
            goal="""Create the most amazing travel itineraries with budget and 
        packing suggestions for the city""",
            backstory="""Specialist in travel planning and logistics with 
        decades of experience""",
            tools=[
                SearchTools.search_internet,
                BrowserTools.scrape_and_summarize_website,
                CalculatorTools.calculate,
            ],
            verbose=True)

@agentops.record_function('main')
def main():
    trip_agents = TripAgents()
    
    city_selection = trip_agents.city_selection_agent()
    local_expert = trip_agents.local_expert()
    travel_concierge = trip_agents.travel_concierge()
    
    # Here you can use these agents in your travel planning logic
    # For example:
    # result = city_selection.run("Find the best city to visit in Europe in July")
    # print(result)

if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        agentops.log_error(str(e))
    finally:
        agentops.end_session('Success')