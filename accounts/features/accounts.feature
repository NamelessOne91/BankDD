Feature: manage bank account
    In order to manage my money
    As a user
    I need to be able to execute operations on my bank account

    Scenario: Open Bank account
        Given I have a new bank account
        When I open it with a starting sum of 1000 euro
        Then it should have a balance of 1000 euro 
        And it should be an active account
        And it should have today as opening date
    
    Scenario: Open Bank Account with amount = 0
        Given I have a new bank account
        When I open it with a starting sum of 0 euro
        Then I should get an invalid amount error
    
    Scenario: Open Bank Account with negative amount
        Given I have a new bank account
        When I open it with a starting sum of -1000 euro
        Then I should get an invalid amount error
    
    Scenario: Open Bank Account which is already active
        Given I have an active bank account
        When I open it with a starting sum of 1000 euro
        Then I should get an error because the account is already open
