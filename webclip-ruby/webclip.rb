# Import of external dependencies
require 'json' # json parser
require 'FileUtils' # file handling

# Import of own classes
require_relative 'system/ProjectManager'


# define function for printing custom error messages
def printe(message)
  puts ""
  puts "=== FEHLER ==="
  puts message
  puts "=============="
  puts ""
  abort
end


# Beign
system 'cls'

ProjectManager.init

# Manually select project
project = ProjectManager.projectSelector
# autoselect ilmes for testing
# project = ProjectManager.getProject('ilmes')


outpath = project.render

system 'cls'
puts '+-------------------------------+'
puts '| Projekt erfolgreich erstellt! |'
puts '+-------------------------------+'
puts ''
puts "Output-Ordner:"
puts File.dirname(__FILE__)+'/'+outpath
puts ""
