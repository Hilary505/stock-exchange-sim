## How to Run

    Build: Navigate to the root of your process-optimizer directory.

        go build -o stock_exchange .

        go build -o checker ./checker

    Run the Optimizer:

        ./stock_exchange examples/simple.txt 1

        ./stock_exchange examples/finite_car_factory.txt 10

        ./stock_exchange examples/infinite_ecosystem.txt 0.5 (Run with a short timeout)

    Run the Checker:

        ./checker examples/simple.txt examples/simple.log

        ./checker examples/finite_car_factory.txt examples/finite_car_factory.log

