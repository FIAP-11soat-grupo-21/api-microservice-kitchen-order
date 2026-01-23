Feature: Kitchen Order Management

  Scenario: Create a new kitchen order
    Given the kitchen order data is valid with order ID "order-123"
    When I send a request to create a new kitchen order
    Then the kitchen order should be created successfully

  Scenario: Find kitchen order by ID
    Given a kitchen order exists with ID "ko-123"
    When I send a request to find the kitchen order by ID
    Then the kitchen order should be returned successfully

  Scenario: Update kitchen order status
    Given a kitchen order exists with ID "ko-456"
    And the new status is "PREPARING"
    When I send a request to update the kitchen order status
    Then the kitchen order status should be updated successfully
