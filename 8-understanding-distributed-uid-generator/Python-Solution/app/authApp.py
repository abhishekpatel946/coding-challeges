from fastapi import FastAPI, Request
from fastapi import status as HttpStatus
from starlette.responses import JSONResponse

from distributed_uid_generator.core import DistributedUIDGenerator
from distributed_uid_generator.optimized_core import (
    SnowflakeGenerator,
    TestSnowflakeGenerator,
)

authApp = FastAPI(
    prefix="/api/v1", docs_url="/docs", openapi_url="/openapi.json", redoc_url="/redoc"
)


@authApp.get("/generate")
async def generate(request: Request):
    request_id = request.state.request_id
    uid = DistributedUIDGenerator().generate_uid()
    return JSONResponse(
        content={
            "responseData": {
                "request_id": request_id,
                "unique_id_binary_format": uid[0],
                "unique_id_number_format": uid[1],
            },
            "message": "Unique ID successfully generated",
            "success": True,
            "code": HttpStatus.HTTP_200_OK,
        },
        status_code=HttpStatus.HTTP_200_OK,
    )


@authApp.get("/optimized-generate")
async def optimize_generate(request: Request):
    request_id = request.state.request_id
    uid = SnowflakeGenerator(datacenter_id=1, machine_id=1).generate_id()
    return JSONResponse(
        content={
            "responseData": {
                "request_id": request_id,
                "unique_id": uid,
            },
            "message": "Unique ID successfully generated",
            "success": True,
            "code": HttpStatus.HTTP_200_OK,
        },
        status_code=HttpStatus.HTTP_200_OK,
    )


@authApp.get("/test/optimized-generate")
async def test_optimized_fn(total_num_ids: int = 10_000):
    total_time_taken = TestSnowflakeGenerator(total_num_ids).test()
    return JSONResponse(
        content={
            "responseData": {
                "total_time_taken": f"Generated {total_num_ids} IDs in {total_time_taken:.3f} seconds"
            },
            "message": "Unique ID successfully generated",
            "success": True,
            "code": HttpStatus.HTTP_200_OK,
        },
        status_code=HttpStatus.HTTP_200_OK,
    )
