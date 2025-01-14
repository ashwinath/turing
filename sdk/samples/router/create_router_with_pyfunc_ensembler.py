import turing
import turing.batch
import turing.batch.config
import turing.router.config.router_config
from turing.router.config.route import Route
from turing.router.config.router_config import RouterConfig
from turing.router.config.router_version import RouterStatus
from turing.router.config.resource_request import ResourceRequest
from turing.router.config.log_config import LogConfig, ResultLoggerType
from turing.router.config.router_ensembler_config import PyfuncRouterEnsemblerConfig
from turing.router.config.experiment_config import ExperimentConfig

from typing import List, Any


# To register a pyfunc ensembler to be used in a Turing router, implement the `turing.ensembler.PyFunc` interface
class SampleEnsembler(turing.ensembler.PyFunc):
    """
    A simple ensembler, that returns the value corresponding to the version that has been specified in the
    `features` in each request. This value if obtained from the route responses found in the `predictions` in each
    request.

    If no version is specified in `features`, return the sum of all the values of all the route responses in
    `predictions` instead.

    e.g. The values in the route responses (`predictions`) corresponding to the versions, `a`, `b` and `c` are 1, 2
         and 3 respectively.

         For a given request, if the version specified in `features` is "a", the ensembler would return the value 1.

         If no version is specified in `features`, the ensembler would return the value 6 (1 + 2 + 3).
    """
    # `initialize` is essentially a method that gets called when an object of your implemented class gets instantiated
    def initialize(self, artifacts: dict):
        pass

    # Each time a Turing Router sends a request to a pyfunc ensembler, ensemble will be called, with the request payload
    # being passed as the `features` argument, and the route responses as the `predictions` argument.
    #
    # If an experiment has been set up, the experiment returned would also be passed as the `treatment_config` argument.
    #
    # The return value of `ensemble` will then be returned as a `json` payload to the Turing router.
    def ensemble(
            self,
            features: dict,
            predictions: List[dict],
            treatment_config: dict) -> Any:
        # Get a mapping between route names and their corresponding responses
        routes_to_response = dict()
        for prediction in predictions:
            routes_to_response[prediction["route"]] = prediction

        if "version" in features:
            return routes_to_response[features["version"]]["data"]["value"]
        else:
            return sum(response["data"]["value"] for response in routes_to_response.values())


def main(turing_api: str, project: str):
    # Initialize Turing client
    turing.set_url(turing_api)
    turing.set_project(project)

    # Register an ensembler with Turing:
    ensembler = turing.PyFuncEnsembler.create(
        name="sample-ensembler-1",
        ensembler_instance=SampleEnsembler(),
        conda_env={
            'dependencies': [
                'python>=3.7.0',
                # other dependencies, if required
            ]
        }
    )
    print("Ensembler created:\n", ensembler)

    # Build a router config in order to create a router
    # Create some routes
    routes = [
        Route(
            id='control',
            endpoint='http://control.endpoints/predict',
            timeout='20ms'
        ),
        Route(
            id='experiment-a',
            endpoint='http://experiment-a.endpoints/predict',
            timeout='20ms'
        )
    ]

    # Create an experiment config (
    experiment_config = ExperimentConfig(
        type="nop"
    )

    # Create a resource request config for the router
    resource_request = ResourceRequest(
        min_replica=0,
        max_replica=2,
        cpu_request="500m",
        memory_request="512Mi"
    )

    # Create a log config for the router
    log_config = LogConfig(
        result_logger_type=ResultLoggerType.NOP
    )

    # Create an ensembler for the router
    ensembler_config = PyfuncRouterEnsemblerConfig(
        project_id=1,
        ensembler_id=1,
        resource_request=ResourceRequest(
            min_replica=0,
            max_replica=2,
            cpu_request="500m",
            memory_request="512Mi"
        ),
        timeout="60ms",
    )

    # Create the RouterConfig instance
    router_config = RouterConfig(
        environment_name="id-dev",
        name="router-with-pyfunc-ensembler",
        routes=routes,
        rules=[],
        default_route_id="test",
        experiment_engine=experiment_config,
        resource_request=resource_request,
        timeout="100ms",
        log_config=log_config,
        ensembler=ensembler_config
    )

    # Create a new router using the RouterConfig object
    new_router = turing.Router.create(router_config)
    print(f"You have created a router with id: {new_router.id}")

    # Wait for the router to get deployed
    try:
        new_router.wait_for_status(RouterStatus.DEPLOYED)
    except TimeoutError:
        raise Exception(f"Turing API is taking too long for router {new_router.id} to get deployed.")

    # 2. List all routers
    routers = turing.Router.list()
    for r in routers:
        print(r)


if __name__ == '__main__':
    import fire
    fire.Fire(main)
