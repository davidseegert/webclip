# -*- coding: UTF-8 -*-

require 'FileUtils'
require_relative 'Project'

class ProjectManager
    @@projects = Hash.new
    @@activeProject = nil

    def self.init
        # loading configuraton
        begin
          contents = File.read('config.json')
        rescue
          puts "=== FEHLER ==="
          puts "Die Datei \"config.json\" im Hauptverzeichnis wurde nicht gefunden."
          puts "1) Stellen Sie sicher, dass die Datei existiert."
          puts "2) Stellen Sie sicher, dass das Script aus dem Hauptverzeichnis heraus gestartet wurde."
          puts "=============="
        end

        begin
          config = JSON.parse(contents)
        rescue
          puts "=== FEHLER ==="
          puts "Die Datei \"config.json\" im Hauptverzeichnis konnte nicht richtig eingelesen werden."
          puts "Stellen Sie sicher, dass die Datei nur valides JSON enthält."
          puts "=============="
        end

        # add all projects to manager
        projectsFolder = Dir.entries(config['projectsPath'])

        projectsFolder.each do |project|
            # check for each folder if its a project-folder
            if project != '.' and project != '..' and File.directory?(File.join(config['projectsPath'],project)) and File.exists?(File.join(config['projectsPath'],project,'config.json'))
                ProjectManager.addProject(File.join(config['projectsPath'],project))
            end
        end

    end

    def self.addProject(path)
        project = Project.new(path)
        # create new projet-object and add it to projects-array
        @@projects[project.getName] = project
    end

    def self.getProject(name)
        @@activeProject = @@projects[name]
        return self.getActiveProject
    end

    def self.getActiveProject
        return @@activeProject
    end

    def self.projectSelector
        @@projects.each do |projectName, project|
            puts projectName
        end

        while true do
          puts "\nEnter name of project to compile:"
          projectName = gets.chomp
          if @@projects.include? projectName
            return self.getProject(projectName)
          else
            puts "\"#{projectName}\" is no project."
          end
        end
    end

end
