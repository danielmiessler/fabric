# This pattern generates questions to help a student to review the main concepts of the learning objectives provided. 
For more clarity the input data should define the subject and the list of learning objectives.

Example input query:

"""
# Optional to be defined here or in the context file
[Student Level: High school student]

Subject: Machine Learning

Learning Objectives:
    * Define machine learning
    * Define supervised learning
    * Define unsupervised learning
    * Define a regression model
"""

# Example run:

Copy the input query to the clipboard.
```bash
xclip -selection clipboard -o | fabric -sp create_quiz
```


## Meta

- **Author**: Marc Andreu (marc@itqualab.com)
- **Version Information**: Marc Andreu's main `create_quiz` version.
- **Published**: May 6, 2024
