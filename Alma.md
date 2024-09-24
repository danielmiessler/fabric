## Document Purpose

This document captures the SPQA policy and State for Alma Security, a security startup out of Redwood City, Ca.

This is part of the SPQA context that will be used to answer questions and create artifacts for the company, e.g., company strategy, security strategy, quarterly security reports (QSRs), project plans, recommendations on which projects to undertake, which investments to take and avoid, and other such decisions.

A major aspect of the SPQA system is the definition of the company's mission, goals, KPIs, and challenges. These shape everything within the company and thus should be used to shape the recommendations made when asked.

In addition to the clearly stated goals and other defining characteristics listed above, there will also be a streaming list of updates coming into this system using the Activity document.

Those will be changes, updates, or modifications to the direction of the company. For example, if Goal number 4 is to build a new datacenter in Boise, Idaho, but we see an update in the Activity section that says we've lost the ability to build in Boise, we should consider goal #4 out of the picture for prioritization and other decision purposes. In other words, the streaming activity log into this document should be considered updates to the core content.

## Company History

Alma Security was started by Chris Meyers, who was previously at Sigma Systems as CTO and HPE as a senior security engineer.

He started the company becuase, "I saw a gap in the authentication market, where companies were only looking at one or two aspects of one's identity to do authentication. They we're looking at the whole picture and turning that into a continuous authentication story."

## Company Mission

The mission of Alma Security is to ensure businesses can continuously authenticate their users using their whole selves.

## Company Goals (G1 means goal 1, G2 is goal 2, etc. Treat each item (goal/kpi/etc) as half as important as the one before it.)

NOTE: Some goals are things like project rollouts which serve the higher goals. In that case they shouldn't always be considered so much lower priority because one is serving the other.

## Company Goals

- G1: Achieve 20% market share by January 2025
- G2: Hit 10000 active customers by January 2025
- G3: Hit a customer trust score of 90+% by January 2025
- G4: Get churn below 5% by August 2024
- G5: Launch in Europe by August 2024
- G6: Launch in India by November 2024
- G7: Launch Mood-monitor integration by February 2024
- G8: Launch partnership with Apple Passkeys by June 2024

## Company KPIs

- K1: Current marketshare percentage
- K2: Number of active customers
- K3: Current churn percentage
- K4: Launched_in_Europe (yes/no)
- K4: Launched_in_India (yes/no)

-----------------------------------------------------------------------------------------------------------------------

## Security Team Mission 

- SM1: Protect Alma Security's customers and intellectual property from security and privacy incidents.

## Security Team Goals

- SG1: Secure all customer data -- especially biometric -- from security and privacy incidents. 
- SG2: Protect Alma Security's intellectual property from being captured by unathorized parties.
- SG3: Reach a time to detect malicious behavior of less than 4 minutes by January 2025
- SG4: Ensure the public trusts our product, because it's an authentication product we can't survive if people don't trust us.
- SG5: Reach a time to remediate critical vulnerabilties on crown jewel systems of less than 16 hours by August 2025
- SG6: Reach a time to remediate critical vulnerabilties on all systems of less than 3 days by August 2025
- SG7: Complete audit of Apple Passkey integration by February 2025
- SG8: Complete remediation of Apple Passkey vulns by February 2025

## Security Team KPIs (How we measure the team)

- SK1: TTD: Time to detect malicious behavior (Minutes)
- SK1: TTI: Time to begin investigation of malicious behavior (Minutes)
- SK3: TTR-CJC: Time to remediate critical vulnerabilities on crown jewel systems (Hours)
- SK3: TTR-C: Time to remediate critical vulnerabilities on all systems (Hours)
- SK4: PT: Public trust score (Complete, Significant, Moderate, Minimal, Distrust, N/A)

## Risk Register (The things we're most worried about)

- R1: Our infrastructure security team is understaffed by 50% after 5 key people left 
- R2: We are not currently monitoring our external perimeter for attack surface related vulnerabilities like open ports, listening applications, unknown hosts, unknown subdomains pointing to these things, etc. We only do scans once every couple of months and we don't really have anyone to look at the results
- R3: It takes us multiple days to investigate potential malicious behavior on our systems.
- R4: We lack a full list of our assets, including externally facing hosts, S3 buckets, etc., which make up our attack surface
- R5: We have a low public trust score due to the events of 2022.

## Security Team Narrative

### Background

Alma hired a new security team starting in January of 2023 and we have been building out the program since then. The philosophy and approach for the security team is to explicitly articulate what we believe the highest risks are to Alma, to deploy targeted strategies to address those risks, and to use clear, transparent KPIs to show progress towards our goals over time.

### Current Risks

So our risk register looks like this:

1. We are understaffed by 50% after 5 key people left in 2022
2. Our perimeter is not being monitored for attack surface related vulnerabilities
3. It takes us too long to detect and start investigating malicious behavior on our systems
4. We do not have a full list of our assets, which makes it difficult to know what we need to protect
5. We have a low public trust score due to the events of 2022 

### Strategies

As such, our strategies are as follows:

1. Hire 5 more A-tier security professionals
2. Purchase and implement an attack surface management solution
3. Invest in our detection and response capabilities
4. Purchase an asset inventory system that integrates with our attack surface management tool
5. Leverage PR to share as much of our progress as possible with the public to rebuild trust

### How We're Doing

We believe being transparent about our progress is key to everything, and for that reason we maintain a limited number of KPIs that we update every quarter. These metrics will not change often. They will remain consistent so that it's easy to track how we're spending our resources and the progress we're making.

Those KPIs are:

1. Time to detect malicious behavior
2. Time to start investigating malicious behavior
3. Time to remediate critical vulnerabilities on crown jewel systems
4. Time to remediate critical vulnerabilities on all systems
5. Our public trust score

As of $DATE$, our KPIs for these are currently:

$GIVE CURRENT KPIs from the Activity section below$

$INSERT GRAPHS OF KPI PROGRESS OVER TIME HERE$

## Security Team Strategies

- STS1: Hire 5 more A-tier security professionals
- STS2: Purchase an attack surface management solution
- STS3: Invest in our detection and response capabilities
- STS4: Purchase an asset inventory system that integrates with our attack surface management tool

## Infrastructure Notes (a basic description of our tech stack and various context around it)

- We currenty have no WAF protecting our main web app but we're considering adding one
- We have had some issues with S3 buckets becoming public, or being set up as public, which has lead to some close calls with customer data almost being exposed.
- alma.amazon-domain.com is our primary S3 bucket that contains everything, but it's not public readable or listable
- We have a root account for our AWS account that doesn't yet have 2FA on it, but we're working on fixing that within a few weeks (but it's been open for a few months)
- We also use Postgres for all our databases.
- Developers have root access to the all kubernetes nodes via SSH on port 45,001 using a shared developer key issued during laptop provisioning. 
- We're a kubernetes shop and do everything through AWS
- We're logging most stuff to Cloudtrail and we kind of use guarduty, but we don't have a 24/7 team to monitor alerts and logs. We should add that to our list of challenges next time we update our overarll policy
- We also have a Windows infrastructure because some key personell came from Microsoft. The DC is hosted in our head office which is in Redwood City, and anyone who works in that office (most of the 300 employees) uses that to log in when tehy start work. The domain is ALMA.
- There's a domain-joined fileserver running Windows 2012 that most people use to upload new ideas and plans for new products. It uses Windows authentication from the domain. 
- We use a palo alto firewall with 2fa using windows authenticator tied to SSO.
- The name of the AI system doing all this context creation using SPQA is Alma, which is also the name of the company.
- We use Workday for HR stuff. Slack for realtime communications. Outlook 365 as a service. Sentinel One on the workstations and laptops. Servers in AWS are mostly Amazon Linux 2 with a few Ubuntu boxes that are a few years old. 
- We also primarily use Postgres for all of our systems.

## Team

TEAM MEMBER | TEAM ASSIGNED | SKILLS | PAY LEVEL | LOCATION | PROJECTS

Nadia Khan | Detection and Response | D&R (Expert), AWS (Strong), Python (Expert), Kubernetes (Basic), Postgres (Basic) | $249K | Redwood City
Chris Magann | Vulnerability Management | VM (Expert), AWS (Strong), Python (Basic), Postgres (Basic) | $212K | Redwood City
Tigan Wang | Vulnerability Management | VM (Expert), AWS (Strong), Python (Basic), Postgres (Basic) | $217K | Redwood City

## Projects

PROJECT NAME | PROJECT DESCRIPTION | PROJECT PRIORITY | PROJECT MEMBERS | START DATE | END DATE | STATUS | PROJECT COST

WAF Install | Install a WAF in front of our main web app | Critical | Nadia Khan | 2024-01-01 - Ongoing | In Progress | $112K one-time, $9K/month

Multi-Factor Authentication (MFA) Rollout | Implement MFA across all internal and external systems | Critical | Chris Magaan | 2024-01-15 | 2024-05-01 | Planned | $80K one-time, $5K/month 

Procure and Implement ASM | Implement continuous monitoring for attack surface vulnerabilities | High | Tigan Wang | 2024-02-15 | 2024-06-15 | Not Started | $75K one-time, $6K/month 

Data Encryption Upgrade | Upgrade encryption protocols for all sensitive data | Medium | Nadia Khan | 2024-04-01 | 2024-08-01 | Planned | $95K one-time 

Incident Response Enhancement | Develop and implement a 24/7 incident response team | High | Nadia Khan | 2024-03-01 | 2024-07-01 | In Progress | $150K one-time, $10K/month 

Cloud Security Optimization | Optimize AWS cloud security configurations and practices | Medium | Tigan Wang | 2024-02-01 | 2024-06-01 | In Progress | $100K one-time, $8K/month 

S3 Bucket Security | Review and secure all S3 buckets to prevent data breaches | High | Chris Magaan | 2024-01-10 | 2024-04-10 | In Progress | $70K one-time, $5K/month

SQL Injection Mitigation | Implement measures to eliminate SQL injection vulnerabilities | High | Tigan Wang | 2024-01-20 | 2024-05-20 | Not Started | $60K one-time 

## CURRENT STATE (KPIs, Metrics, Project Activity Updates, etc.)
- October 2022: Current time to detect malicious behavior is 81 hours
- October 2022: Current time to start investigating malicious behavior is 82 hours
- October 2022: Current time to remediate critical vulnerabilities on crown jewel systems is 21 days
- October 2022: Current time to remediate critical vulnerabilities on all systems is 51 days
- January 2023: Current time to detect malicious behavior is 62 hours
- January 2023: Current time to start investigating malicious behavior is 72 hours
- January 2023: Current time to remediate critical vulnerabilities on crown jewel systems is 17 days
- January 2023: Current time to remediate critical vulnerabilities on all systems is 43 days
- July 2023: Current time to detect malicious behavior is 29 hours
- July 2023: Current time to start investigating malicious behavior is 41 hours
- July 2023: Current time to remediate critical vulnerabilities on crown jewel systems is 12 days
- July 2023: Current time to remediate critical vulnerabilities on all systems is 29 days
- November 2023: Current time to start detect malicious behavior is 12 hours
- November 2023: Current time to start investigating malicious behavior is 16 hours
- November 2023: Current time to remediate critical vulnerabilities on crown jewel systems is 9 days
- November 2023: Current time to remediate critical vulnerabilities on all systems is 17 days
- February 2024: Started attack surface management vendor selection process
- January 2024: Current time to start detect malicious behavior is 9 hours
- January 2024: Current time to start investigating malicious behavior is 14 hours
- January 2024: Current time to remediate critical vulnerabilities on crown jewel systems is 8 days
- January 2024: Current time to remediate critical vulnerabilities on all systems is 12 days
- March 2024: We're now remediating crits on crown jewels in less than 6 days
- April 2024: We're now remediating all criticals within 11 days
- July 2024: Criticals are now being fixed in 9 days
- On August 5 we got remediation of critical vulnerabilities down to 7 days
