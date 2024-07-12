from crewai import Task
from textwrap import dedent
from datetime import date
import agentops
import os

# Initialize AgentOps
AGENTOPS_API_KEY = os.getenv("AGENTOPS_API_KEY")
if not AGENTOPS_API_KEY:
    raise ValueError("AGENTOPS_API_KEY not found in environment variables")

agentops.init(AGENTOPS_API_KEY)

@agentops.track_agent(name='TripTasks')
class TripTasks:

    @agentops.record_function('identify_task')
    def identify_task(self, agent, origin, cities, interests, range):
        return Task(description=dedent(f"""
            Analyze and select the best city for the trip based 
            on specific criteria such as weather patterns, seasonal
            events, and travel costs. This task involves comparing
            multiple cities, considering factors like current weather
            conditions, upcoming cultural or seasonal events, and
            overall travel expenses. 
            
            Your final answer must be a detailed
            report on the chosen city, and everything you found out
            about it, including the actual flight costs, weather 
            forecast and attractions.
            {self.__tip_section()}

            Traveling from: {origin}
            City Options: {cities}
            Trip Date: {range}
            Traveler Interests: {interests}
        """),
                    agent=agent)

    @agentops.record_function('gather_task')
    def gather_task(self, agent, origin, interests, range):
        return Task(description=dedent(f"""
            As a local expert on this city you must compile an 
            in-depth guide for someone traveling there and wanting 
            to have THE BEST trip ever!
            Gather information about  key attractions, local customs,
            special events, and daily activity recommendations.
            Find the best spots to go to, the kind of place only a
            local would know.
            This guide should provide a thorough overview of what 
            the city has to offer, including hidden gems, cultural
            hotspots, must-visit landmarks, weather forecasts, and
            high level costs.
            
            The final answer must be a comprehensive city guide, 
            rich in cultural insights and practical tips, 
            tailored to enhance the travel experience.
            {self.__tip_section()}

            Trip Date: {range}
            Traveling from: {origin}
            Traveler Interests: {interests}
        """),
                    agent=agent)

    @agentops.record_function('plan_task')
    def plan_task(self, agent, origin, interests, range):
        return Task(description=dedent(f"""
            Expand this guide into a a full 7-day travel 
            itinerary with detailed per-day plans, including 
            weather forecasts, places to eat, packing suggestions, 
            and a budget breakdown.
            
            You MUST suggest actual places to visit, actual hotels 
            to stay and actual restaurants to go to.
            
            This itinerary should cover all aspects of the trip, 
            from arrival to departure, integrating the city guide
            information with practical travel logistics.
            
            Your final answer MUST be a complete expanded travel plan,
            formatted as markdown, encompassing a daily schedule,
            anticipated weather conditions, recommended clothing and
            items to pack, and a detailed budget, ensuring THE BEST
            TRIP EVER, Be specific and give it a reason why you picked
            # up each place, what make them special! {self.__tip_section()}

            Trip Date: {range}
            Traveling from: {origin}
            Traveler Interests: {interests}
        """),
                    agent=agent)

    @agentops.record_function('__tip_section')
    def __tip_section(self):
        return "If you do your BEST WORK, I'll tip you $100!"

@agentops.record_function('main')
def main():
    # This function can be used to run any initialization or main logic
    trip_tasks = TripTasks()
    
    # Example usage:
    # agent = SomeAgent()  # You would need to import and create an appropriate agent
    # identify_task = trip_tasks.identify_task(agent, "New York", "Paris, London, Rome", "museums, food", "June 1-7, 2024")
    # gather_task = trip_tasks.gather_task(agent, "New York", "museums, food", "June 1-7, 2024")
    # plan_task = trip_tasks.plan_task(agent, "New York", "museums, food", "June 1-7, 2024")
    
    # Here you can use these tasks in your trip planning logic

if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        agentops.log_error(str(e))
    finally:
        agentops.end_session('Success')