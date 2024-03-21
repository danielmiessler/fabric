## What this Pull Request (PR) does
PR includes two new patterns:

1) `create_report_finding`
    * This takes either a file or text via echo for example and creates a pentest report finding that includes the following sections:
        * title, description, risk, remedation, external references (please check these), one-sentence-summary, quotes.
    * example usage: echo "Username Enumeration: Forgotten Password Functionality: The application returns if an account exists or not, which allows an attacker to enumerate valid user accounts via email address" | create_report_finding
    
2) `improve_report_finding`
    * This takes either a file or text via echo for example and creates an improved pentest report finding that includes the following sections:
        * title, description, risk, remedation, external references (please check these), one-sentence-summary, quotes.
    * example usage: cat sanitised_report_finding.txt (should have title, description, remediation sections) | improve_report_finding

Additionally, this PR includes a Github helper script for automating the Github contributing workflow. This allows you to:

1) Update your fork with the main repo to ensure you're working on a current version.
2) Create a new branch.
3) Push changes to your branch (or new branch).
4) Create a PR using a markdown file to populate the body.

## Example Output from `create_report_finding`:
### Username Enumeration: Forgotten Password Functionality

#### Description
The application in question has a security flaw within its forgotten password functionality. Specifically, when a user attempts to reset their password using an email address, the application responds differently depending on whether the email address is associated with an existing account. This behavior inadvertently provides attackers with a means to confirm the existence of valid user accounts. By systematically submitting various email addresses through this functionality, an attacker can compile a list of valid accounts for further malicious activities, such as targeted phishing attacks or brute force password attempts.

#### Risk
This vulnerability poses a significant risk as it directly compromises user privacy and security. The ability for an attacker to enumerate valid user accounts elevates the risk of targeted attacks. Users with identified accounts may become victims of phishing campaigns designed to extract more sensitive information or deceive them into compromising their account security. Furthermore, knowing which accounts are valid can aid an attacker in focusing their efforts on existing accounts when attempting password breaches, making the attack more efficient and likely to succeed.

#### Recommendations
- Implement a uniform response message for all password reset attempts, regardless of whether the email address is associated with an existing account or not.
- Employ CAPTCHA mechanisms to prevent automated scripts from performing mass enumeration attempts.
- Rate limit the number of password reset requests that can be made from a single IP address within a given timeframe to deter enumeration attacks.
- Monitor and log all password reset attempts to detect and respond to potential enumeration activities.
- Educate users on the importance of using unique, strong passwords for their accounts to mitigate the risk of unauthorized access should their email address be enumerated.
- Consider implementing multi-factor authentication (MFA) as an additional layer of security for account access, reducing the impact of account enumeration.

#### References
- [OWASP Guide to Authentication](https://owasp.org/www-project-cheat-sheets/cheatsheets/Authentication_Cheat_Sheet.html)
- [NIST Recommendations on Digital Identity Guidelines](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-63b.pdf)
- [CWE-203: Information Exposure Through Discrepancy](https://cwe.mitre.org/data/definitions/203.html)

#### One-Sentence-Summary:
The forgotten password functionality reveals if an email is linked to an account, enabling attackers to identify valid user accounts.

#### Trends:
- Increasing sophistication of automated scripts used by attackers for account enumeration.
- Growing awareness and adoption of multi-factor authentication (MFA) as a countermeasure.
- Enhanced focus on privacy regulations prompting better security practices around user data.
- Rise in targeted phishing attacks leveraging enumerated account information.
- Shift towards uniform error responses across web applications to mitigate enumeration risks.

#### Quotes:
- "The application responds differently depending on whether the email address is associated with an existing account."
- "This behavior inadvertently provides attackers with a means to confirm the existence of valid user accounts."
- "Users with identified accounts may become victims of phishing campaigns."
- "Knowing which accounts are valid can aid an attacker in focusing their efforts."
- "Implement a uniform response message for all password reset attempts."
- "Employ CAPTCHA mechanisms to prevent automated scripts."
- "Rate limit the number of password reset requests from a single IP address."
- "Educate users on the importance of using unique, strong passwords."
- "Consider implementing multi-factor authentication (MFA) as an additional layer of security."
- "Increasing sophistication of automated scripts used by attackers for account enumeration."%          
