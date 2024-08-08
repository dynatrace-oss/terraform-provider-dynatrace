resource "dynatrace_openpipeline" "this" {
  custom_base_path = "i dont know"
  endpoints {
    endpoint {
       
        base_path = "still dont know"
        processors {
            processor{
              dql_processor{

              }
            }

            processor{
              dql_processor {

              }
            }

            processor {
              fields_add_processor {

              }
            }

            
           
            
   
        }
    } 
  }