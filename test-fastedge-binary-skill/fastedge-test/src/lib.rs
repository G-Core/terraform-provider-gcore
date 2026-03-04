use fastedge::{http_handler, Request, Response};

#[http_handler]
fn handle_request(_req: Request) -> Response {
    Response::new(200, "Hello from FastEdge!")
}
