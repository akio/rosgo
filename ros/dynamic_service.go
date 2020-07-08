package ros

// IMPORT REQUIRED PACKAGES.

import (
	"strings"

	"github.com/edwinhayes/rosgo/libgengo"
	"github.com/pkg/errors"
)

// DEFINE PUBLIC STRUCTURES.

// DynamicServiceType abstracts the schema of a ROS Service whose schema is only known at runtime.
// DynamicServiceType implements the rosgo ServiceType interface, allowing it to be used throughout rosgo in the same manner as message schemas generated
// at compiletime by gengo.
type DynamicServiceType struct {
	name    string
	md5sum  string
	text    string
	reqType MessageType
	resType MessageType
}

// DynamicService abstracts an instance of a ROS Service whose type is only known at runtime.  The schema of the message is denoted by the referenced DynamicServiceType, while the
// Request and Response references the rosgo Messages it implements.  DynamicService implements the rosgo Service interface, allowing
// it to be used throughout rosgo in the same manner as service types generated at compiletime by gengo.

type DynamicService struct {
	dynamicType *DynamicServiceType
	Request     Message
	Response    Message
}

// DEFINE PRIVATE STRUCTURES.

// DEFINE PUBLIC GLOBALS.

// DEFINE PRIVATE GLOBALS.

// DEFINE PUBLIC STATIC FUNCTIONS.

// NewDynamicServiceType generates a DynamicServiceType corresponding to the specified typeName from the available ROS service definitions; typeName should be a fully-qualified
// ROS service type name.  The first time the function is run, a message/service 'context' is created by searching through the available ROS definitions, then the ROS service to
// be used for the definition is looked up by name.  On subsequent calls, the ROS service type is looked up directly from the existing context.

func NewDynamicServiceType(typeName string) (*DynamicServiceType, error) {
	return newDynamicServiceTypeNested(typeName, "")
}

// newDynamicServiceTypeNested generates a DynamicServiceType from the available ROS definitions.  The first time the function is run, a message/service 'context' is created by
// searching through the available ROS definitions, then the ROS service type to use for the defintion is looked up by name.  On subsequent calls, the ROS service type
// is looked up directly from the existing context.  This 'nested' version of the function is able to be called recursively, where packageName should be the typeName of the
// parent ROS services; this is used internally for handling complex ROS services.
func newDynamicServiceTypeNested(typeName string, packageName string) (*DynamicServiceType, error) {
	// Create an empty message type.
	m := new(DynamicServiceType)

	// Create a message context if for some reason it does not exist yet, as it also contains service definitions
	if context == nil {
		// Create context for our ROS install.
		c, err := libgengo.NewMsgContext(strings.Split(GetRuntimePackagePath(), ":"))
		if err != nil {
			return nil, err
		}
		context = c
	}
	// We need to try to look up the full name, in case we've just been given a short name.
	fullname := typeName

	_, ok := context.GetMsgs()[fullname]
	if !ok {
		// Seems like the package_name we were give wasn't the full name.

		// Messages in the same package are allowed to use relative names, so try prefixing the package.
		if packageName != "" {
			fullname = packageName + "/" + fullname
		}
	}

	// Load context for the target message.
	spec, err := context.LoadSrv(fullname)
	if err != nil {
		return nil, err
	}

	// Now we know all about the service!
	m.name = spec.ShortName
	m.md5sum = spec.MD5Sum
	m.text = spec.Text
	m.reqType, err = NewDynamicMessageType(spec.Request.FullName)
	if err != nil {
		return nil, errors.Wrap(err, "error generating request type")
	}
	m.resType, err = NewDynamicMessageType(spec.Response.FullName)
	if err != nil {
		return nil, errors.Wrap(err, "error generating request type")
	}

	// We've successfully made a new service type matching the requested ROS type.
	return m, nil

}

// DEFINE PUBLIC RECEIVER FUNCTIONS.

//	DynamicMessageType

// MD5Sum returns the ROS compatible MD5 sum of the service type; required for ros.ServiceType.
func (t *DynamicServiceType) MD5Sum() string {
	return t.md5sum
}

// Name returns the full ROS name of the service type; required for ros.ServiceType.
func (t *DynamicServiceType) Name() string {
	return t.name
}

// NewService creates a new DynamicService instantiating the service type; required for ros.ServiceType.
func (t *DynamicServiceType) NewService() Service {
	// Don't instantiate services for incomplete types.
	if t == nil {
		return nil
	}
	// But otherwise, make a new one.
	d := new(DynamicService)
	d.dynamicType = t
	d.Request = t.RequestType().NewMessage()
	d.Response = t.ResponseType().NewMessage()
	return d
}

// RequestType returns the MessageType of the request message of DynamicServiceType; required for ros.ServiceType.
func (t *DynamicServiceType) RequestType() MessageType {
	return t.reqType
}

// ResponseType returns the MessageType of the response message of DynamicServiceType; required for ros.ServiceType.
func (t *DynamicServiceType) ResponseType() MessageType {
	return t.resType
}

//	DynamicService

// ReqMessage returns the request message of the DynamicService; required for ros.Service.
func (s *DynamicService) ReqMessage() Message {
	return s.Request
}

// ResMessage returns the response message of the DynamicService; required for ros.Service.
func (s *DynamicService) ResMessage() Message {
	return s.Response
}

// Type returns the ROS type of a dynamic service; not required for ros.Service.
func (s *DynamicService) Type() ServiceType {
	return s.dynamicType
}

// DEFINE PRIVATE STATIC FUNCTIONS.

// DEFINE PRIVATE RECEIVER FUNCTIONS.

// ALL DONE.
