import * as React from "react"
import { ListGroup, Table, Button, Badge } from "react-bootstrap"
import * as moment from "moment"
import { AggregatedError, Error, HttpContext, Headers, StoreState, ErrorInstance } from "data/types"
import { ButtonGroup } from "react-bootstrap"
import { setCurrentExceptionIndex, resolveError } from "data/errors"
import { bindActionCreators, Dispatch, AnyAction } from "redux"
import { connect } from "react-redux"

interface ConnectedProps {
  activeError: AggregatedError,
  activeService: string,
  latestExceptionIndex: number,
}

interface DispatchProps {
  setCurrentExceptionIndex: (num: number) => void,
  resolveError: (service: string, errorKey: string) => void
}

type Props = ConnectedProps & DispatchProps


const ErrorComponent: React.FC<Props> = (props) => {

  const calculateNewIndex = (index: number, inc: number, size: number) => {
    return (index + inc % size + size) % size
  }

  const showPreviousException = () => {
    props.setCurrentExceptionIndex(calculateNewIndex(props.latestExceptionIndex, -1, props.activeError.latest_errors.length))
  }

  const showNextException = () => {
    props.setCurrentExceptionIndex(calculateNewIndex(props.latestExceptionIndex, 1, props.activeError.latest_errors.length))
  }

  const renderErrorInstance = (errorInstance: ErrorInstance) => {
    if (!errorInstance) return ""

    return (
      <div>      
        <ListGroup.Item>
          <h4 className="list-group-item-heading"> Class</h4>
          { errorInstance.class }
        </ListGroup.Item>
        { renderMessage(errorInstance.message) }
        { renderStackTrace(errorInstance.stacktrace) }
        { renderCause(errorInstance.cause) }
      </div>
    )
  }

  const renderError = (error: Error) => {
    if (!error) return ""

    return (
      <ListGroup variant="flush">
        <ListGroup.Item>
          <h4 className="list-group-item-heading"> Occurred at</h4>
          {renderTimeOccurrence(error.timestamp)}
        </ListGroup.Item>  
        { renderCurl(error.http_context) }
        { renderErrorInstance(error.error) }
        { renderHttpContext(error.http_context) }
      </ListGroup>
    )
  }

  const renderStackTrace = (stackTrace: string[]) => {
    if (stackTrace === null || stackTrace.length === 0) {
      return ""
    }
    const trace =  stackTrace.map((line) => line + "\n")
    return (
      <ListGroup.Item>
        <h4 className="list-group-item-heading"> Trace</h4>
        <pre className="pre-scrollable"><code>{ trace }</code></pre>
    </ListGroup.Item>
    )
  }

  const renderTimeOccurrence = (ts: number) => {
    return moment(new Date(ts * 1000)).fromNow()
  }

  const renderCause = (cause: ErrorInstance) => {
    if (!cause) {
      return ""
    }

    return (
      <ListGroup.Item>
        <h4 className="list-group-item-heading"> Cause</h4>
        { renderErrorInstance(cause) }
      </ListGroup.Item>
    )
  }

  const renderMessage = (message: string) => {
    if (message === null || message.trim().length === 0) {
      return ""
    }

    return (
      <ListGroup.Item>
        <h4 className="list-group-item-heading"> Message</h4>
        { message }
      </ListGroup.Item>
    )
  }

  const renderCurl = (context: HttpContext) => {
    if (context == null) {
      return ""
    }

    let headers: Headers = context.request_headers == null ? {} : context.request_headers

    let headersString: string = Object.keys(headers).reduce(function(headersString, key) {
      return `${headersString}-H "${key}: ${headers[key]}" `
    }, "")

    return (
      <ListGroup.Item>
        <h4 className="list-group-item-heading"> Curl</h4>
        <pre>curl -X { context.request_method } {headersString} {context.request_url}</pre>
      </ListGroup.Item>
    )
  }

  const renderContextHeadersRow = (key: string, value: string) => {
    return (
      <tr key={key}>
        <td>{`${key}`}</td>
        <td>{`${value}`}</td>
      </tr>
    )
  }

  const renderContextHeaders = (headers: Headers) => {
    return (
      <ListGroup.Item>
        <h4 className="list-group-item-heading"> Headers</h4>
        <Table striped>
          <tbody>
            {Object.keys(headers).map((key) => {
              return renderContextHeadersRow(
                key,
                headers[key]
              );
            })}
          </tbody>
        </Table>
      </ListGroup.Item>
    );
  };

  const renderRequestBody = (request_body: string) => {
    return(
      <ListGroup.Item>
        <h4 className="list-group-item-heading"> Body</h4>
        {request_body}
      </ListGroup.Item>
    );
  };

  const renderHttpContext = (context: HttpContext) => {
    if (context == null) {
      return ""
    }

    return (
      <ListGroup.Item>
        <h4 className="list-group-item-heading"> HTTP Context</h4>
        <ListGroup>
          <ListGroup.Item>
            <h4 className="list-group-item-heading"> Url</h4>
            {context.request_url}
          </ListGroup.Item>
          <ListGroup.Item>
            <h4 className="list-group-item-heading"> Method</h4>
            {context.request_method}
          </ListGroup.Item>
          {context.request_headers ? renderContextHeaders(context.request_headers) : null}
          {context.request_body ? renderRequestBody(context.request_body): null}
        </ListGroup>
      </ListGroup.Item>
    )
  }

  const resolveError = () => {
    props.resolveError(props.activeService, props.activeError.aggregation_key)
  }
    
  const renderAggregatedError = () => {
    return (
      <div className={"grid-component"}>
        <ButtonGroup className="float-right">
          <Button variant="outline-danger" size="sm" onClick={() => resolveError()}>Resolve</Button>
        </ButtonGroup>
        <h3 className="list-group-item-heading"> Summary</h3>
        <ListGroup>
          <ListGroup.Item>
            <h4 className="list-group-item-heading"> Key</h4>
            {props.activeError.aggregation_key}
          </ListGroup.Item>
          <ListGroup.Item>
            <h4 className="list-group-item-heading"> Count</h4>
            {props.activeError.total_count}
          </ListGroup.Item>
          <ListGroup.Item>
            <h4 className="list-group-item-heading"> Severity</h4>
            {props.activeError.severity}
          </ListGroup.Item>
          <ListGroup.Item>
            <h4 className="list-group-item-heading"> First Occurrence</h4>
            {renderTimeOccurrence(props.activeError.created_at)}
          </ListGroup.Item>
        </ListGroup>
        <br/>
        <ButtonGroup className="float-right">
          <Button variant="outline-dark" size="sm" onClick={() => showPreviousException()}>Previous</Button>
          <Button variant="outline-dark" size="sm" onClick={() => showNextException()} >Next</Button>
        </ButtonGroup>
        <h3 className="list-group-item-heading"> Latest Occurrences <Badge variant="light">{props.latestExceptionIndex+1 + "/" + props.activeError.latest_errors.length}</Badge></h3>
        {renderError(props.activeError.latest_errors[props.latestExceptionIndex]) }
      </div>
    )
  }

  return (
    <div>
      { renderAggregatedError() }
    </div>
  )
}

const mapStateToProps = (state: StoreState) => {
  return {
    activeError: state.errorsReducer.activeError,
    activeService: state.errorsReducer.activeService,
    latestExceptionIndex: state.errorsReducer.latestExceptionIndex,
  }
}

const matchDispatchToProps = (dispatch: Dispatch<AnyAction>): DispatchProps => {
  return bindActionCreators({ setCurrentExceptionIndex, resolveError }, dispatch)
}

export default connect(mapStateToProps, matchDispatchToProps)(ErrorComponent)
